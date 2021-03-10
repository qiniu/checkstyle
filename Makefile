all:
	bash -c "mkdir -p temp; cd temp && export GOPATH=`pwd`/temp && go get github.com/tomcz/checkstyle/gocheckstyle"
	temp/bin/gocheckstyle -config=.gostyle .
	@rm -fr temp
