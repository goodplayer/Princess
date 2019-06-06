package controller2

import (
	"net/http"
	"time"

	"github.com/goodplayer/Princess/repo"
	"github.com/goodplayer/Princess/utils/logging"

	"github.com/gin-gonic/gin"
)

func InitRecording(r *gin.Engine) {
	r.GET("/recording", RecordingIndex)
	r.GET("/recording/", RecordingIndex)

	r.POST("/recording/new", RecordingNew)
}

var (
	recordingLogger = logging.NewLogger("recordding")
)

type Recording struct {
	Id               int64
	ArticleId        int64
	SeqId            int64
	Title            string
	Summary          *string
	Content          string
	Attachment       *string
	Config           *string
	Status           int64
	Type             int64
	CreateTime       int64
	ModifyTime       int64
	AttachmentUnikey *string

	CreateTimeString string
	ModifyTimeString string
}

type RecordingType struct {
	Id         int64
	Name       string
	Title      string
	RegEx      *string
	Config     *string
	CreateTime int64
	ModifyTime int64
}

func RecordingNew(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.HTML(http.StatusOK, "500.html", struct {
			Message string
		}{
			Message: "报错了...",
		})
		recordingLogger.Error("parse new recording multipart form error:", err)
		return
	}

	titles, ok := form.Value["title"]
	contents, ok := form.Value["content"]
	recordingUrls, ok := form.Value["recording_url"]
	recordingTypes, ok := form.Value["recording_type"]

	fileHeaders, ok := form.File["attachment"]
	fileHeaders[0].Filename
	f, err := fileHeaders[0].Open()

	form.RemoveAll()

	//TODO
}

func RecordingIndex(c *gin.Context) {

	types, err := queryRecordingTypes()
	if err != nil {
		c.HTML(http.StatusOK, "500.html", struct {
			Message string
		}{
			Message: "报错了...",
		})
		return
	}

	recordings, err := queryRecordings()
	if err != nil {
		c.HTML(http.StatusOK, "500.html", struct {
			Message string
		}{
			Message: "报错了...",
		})
		return
	}

	c.HTML(http.StatusOK, "recording_index.html", struct {
		Types      []RecordingType
		Recordings []Recording
	}{
		Types:      types,
		Recordings: recordings,
	})
}

func queryRecordings() ([]Recording, error) {
	recordingRows, err := repo.Run().Query(`
SELECT distinct on (article_id) id, article_id, seq_id, title, summary, content, unikey, attachment, config, status, type, create_time, modify_time
	FROM public.recording order by article_id, modify_time DESC LIMIT $1;
	`, 20)

	if err != nil {
		recordingLogger.Error("query recordings error:", err)
		return nil, err
	}
	defer func() {
		var _ = recordingRows.Close()
	}()

	var recordings []Recording
	for recordingRows.Next() {
		r := Recording{}
		err := recordingRows.Scan(&r.Id, &r.ArticleId, &r.SeqId, &r.Title, &r.Summary, &r.Content, &r.AttachmentUnikey,
			&r.Attachment, &r.Config, &r.Status, &r.Type,
			&r.CreateTime, &r.ModifyTime)
		if err != nil {
			recordingLogger.Error("foreach recording error:", err)
			return nil, err
		}
		modifyTime := time.Unix(r.ModifyTime/1000, 0)
		r.ModifyTimeString = modifyTime.String()
		createTime := time.Unix(r.CreateTime/1000, 0)
		r.CreateTimeString = createTime.String()
		recordings = append(recordings, r)
	}

	return recordings, nil
}

func queryRecordingTypes() ([]RecordingType, error) {
	typeRows, err := repo.Run().Query(`
SELECT id, name, regex, config, title
	FROM public.recording_type;
	`)

	if err != nil {
		recordingLogger.Error("query recording type error:", err)
		return nil, err
	}
	defer func() {
		var _ = typeRows.Close()
	}()

	var types []RecordingType
	for typeRows.Next() {
		t := RecordingType{}
		err := typeRows.Scan(&t.Id, &t.Name, &t.RegEx, &t.Config, &t.Title)
		if err != nil {
			recordingLogger.Error("foreach recording types error:", err)
			return nil, err
		}
		types = append(types, t)
	}

	return types, nil
}
