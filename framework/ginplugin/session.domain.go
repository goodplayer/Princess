package ginplugin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const defaultSessionKey = "GSESSIONID"

var sessionDomainLogger = logrus.New()

//================= SessionDomainService =======================

func NewSessionDomainService(cookieName string, maxAge int, sessionStore SessionStore) *SessionDomainService {
	return &SessionDomainService{
		CookieName:         cookieName,
		SessionStore:       sessionStore,
		TTL:                maxAge,
		ClientCookieMaxAge: maxAge,
	}
}

type SessionDomainService struct {
	CookieName   string
	SessionStore SessionStore
	TTL          int

	ClientCookieMaxAge int
}

func (s *SessionDomainService) initSessionManager(ctx *gin.Context) {
	// 1. check cookie existing
	cookie, err := ctx.Request.Cookie(s.CookieName)
	var session *sessionItem
	if err == http.ErrNoCookie {
		// 2.1. create a new session if no cookie existing
		// create new session
		session = s.newSessionItem(ctx)
		// prepare cookie parameters
		s.prepareNewSession(session)
	} else {
		// 2.2. process existing cookie
		// prepare session from cookie
		session = s.newSessionItemFromCookie(ctx, cookie)
		// touch existing cookie
		// especially for the case that in the business code flow, the session is not touched
		s.touchCookie(session)
	}
	// 3. prepare session in the context for business use
	ctx.Set(defaultSessionManagerKey, session)
}

func (s *SessionDomainService) prepareNewSession(session *sessionItem) {
	session.prepareCookie(s.CookieName, s.ClientCookieMaxAge, "/", "", true, true)
}

func (s *SessionDomainService) commitSession(session *sessionItem) {
	// Prerequisite assumption:
	// 0. all session data have been recorded in the backed storage
	// Workflow:
	// 1. Generate cookie and prepare the MaxAge inside cookie
	// 2. Attach the cookie in the http response
	//
	// Here we trigger cookie update in http response only once
	// Because multiple set cookie in http response will cause duplicated cookie records with the same name.
	// This may cause some issues when retrieving data from cookie
	http.SetCookie(session.ctx.Writer, session.cookie(s.ClientCookieMaxAge))
}

func (s *SessionDomainService) touchCookie(session *sessionItem) {
	if err := s.SessionStore.TouchIfExists(session.sessionId, s.TTL); err != nil {
		panic(err)
	}
}

func (s *SessionDomainService) newSessionItem(ctx *gin.Context) *sessionItem {
	return &sessionItem{
		data:      map[string][]byte{},
		sessionId: NewSessionId(),
		s:         s,
		ctx:       ctx,
	}
}

func (s *SessionDomainService) newSessionItemFromCookie(ctx *gin.Context, cookie *http.Cookie) *sessionItem {
	session := &sessionItem{
		s:   s,
		ctx: ctx,
		c:   cookie,
	}
	session.recoverSessionItem()
	return session
}

//================= sessionItem =======================

func session(ctx *gin.Context) *sessionItem {
	return ctx.MustGet(defaultSessionManagerKey).(*sessionItem)
}

type sessionItem struct {
	status int // 0 - nothing, 1 - modified, 2 - freezed
	data   SessionData

	sessionId string

	s   *SessionDomainService
	ctx *gin.Context
	c   *http.Cookie
}

func (s *sessionItem) Set(key string, value []byte) {
	if s.status > 1 {
		panic(errors.New("session has been committed"))
	}
	s.data[key] = value
	s.status = 1
}

func (s *sessionItem) Get(key string) []byte {
	return s.data[key]
}

func (s *sessionItem) Delete(key string) {
	if s.status > 1 {
		panic(errors.New("session has been committed"))
	}
	delete(s.data, key)
	s.status = 1
}

func (s *sessionItem) SaveAndFreeze() error {
	if s.status < 1 {
		return nil
	}
	if s.status >= 2 {
		panic(errors.New("session has been committed"))
	}

	if err := s.s.SessionStore.Save(s.sessionId, s.data, s.s.TTL); err != nil {
		return err
	} else {
		s.status = 2
		s.s.commitSession(s)
	}
	return nil
}

func (s *sessionItem) prepareCookie(cookieName string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    fmt.Sprintf("%s=%s", defaultSessionKey, s.sessionId),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteDefaultMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	s.c = cookie
}

func (s *sessionItem) cookie(maxAge int) *http.Cookie {
	if s.c.MaxAge > 0 && (maxAge-s.c.MaxAge < 3600*24*7) { // if the cookie expires in 7 days, then refresh the session key
		s.c.MaxAge = maxAge
	}
	return s.c
}

func (s *sessionItem) recoverSessionItem() {
	valPair := strings.Split(s.c.Value, "=")
	if len(valPair) == 2 {
		if data, err := s.s.SessionStore.Load(valPair[1]); err != nil {
			panic(err)
		} else {
			if data == nil {
				// no data means session not existing in the store, have to initialize a new one
				s.sessionId = NewSessionId()
				s.data = map[string][]byte{}
				s.s.prepareNewSession(s)
			} else {
				s.sessionId = valPair[1]
				s.data = data
			}
		}
	} else {
		// invalid session, but we generate a new session id as if the session is new
		s.sessionId = NewSessionId()
		s.data = map[string][]byte{}
		s.s.prepareNewSession(s)
	}
}
