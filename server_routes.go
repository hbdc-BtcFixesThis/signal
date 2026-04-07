package main

import (
	"io/fs"
	"net/http"
)

func (ss *SignalServer) setHandlers() {
	fs, _ := fs.Sub(static, string(ss.sc.UiDir(nil)))
	ck := ss.CheckAPIKey
	jw := JSONResponseHeadersWrapper
	handlerAuthFunc := func(x func(http.ResponseWriter, *http.Request)) http.Handler {
		return jw(ck(http.HandlerFunc(x)))
	}
	handlerFunc := func(x func(http.ResponseWriter, *http.Request)) http.Handler {
		return jw(http.HandlerFunc(x))
	}

	ss.serveMux.Handle("/", http.FileServer(http.FS(fs)))
	ss.infoLog.Printf("Serving embedded file directory: %s", ss.sc.UiDir(nil))

	// authenticated apis (settings)
	ss.serveMux.Handle("/verify/token", handlerAuthFunc(ss.verifyHandler))
	// jw(ck(http.HandlerFunc(ss.verifyHandler))))
	ss.serveMux.Handle("/max/db/size", handlerAuthFunc(ss.getNodeMaxStorageSize))

	// public
	ss.serveMux.Handle("/new/record", handlerFunc(ss.newRecordAndOrSignal))
	ss.serveMux.Handle("/new/signal", handlerFunc(ss.newRecordAndOrSignal))
	ss.serveMux.Handle("/get/page", handlerFunc(ss.getPage))
	ss.serveMux.Handle("/record/value", handlerFunc(ss.getRecordValue))
	ss.serveMux.Handle("/record/signals", handlerFunc(ss.getRecordSignals))
	ss.serveMux.Handle("/message/template", handlerFunc(ss.getMessageTemplate))

	// server conf
	ss.serveMux.Handle("/sc/TlsCrtFname", handlerFunc(ss.getTlsCrtFname))
	ss.serveMux.Handle("/sc/TlsKeyFname", handlerFunc(ss.getTlsKeyFname))
	ss.serveMux.Handle("/sc/TlsHosts", handlerFunc(ss.getTlsHosts))
	ss.serveMux.Handle("/sc/SignalDataDBFullPath", handlerFunc(ss.getSignalDataDBFullPath))
	//ss.serveMux.Handle("/sc/NostrRelays", jw(http.HandlerFunc(ss.getNostrRelays)))
}
