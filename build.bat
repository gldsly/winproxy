rsrc -manifest uac.manifest -o uac.syso -ico icon\m2.ico
go build -ldflags "-s -w -H windowsgui"

