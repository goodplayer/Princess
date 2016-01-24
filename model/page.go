package model

var _pageUtilBean = pageUtil{}

type pageUtil struct {
}

func PageUtil() pageUtil {
	return _pageUtilBean
}

func (pageUtil) NewPage() *Page {
	page := new(Page)
	//TODO
	return page
}

func (pageUtil) GetNavBar() []*Page {
	//TODO
	return nil
}

type Page struct {
	//TODO
}
