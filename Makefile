
.PHONY: clean tidy

ifeq ($(OS),Windows_NT)
define zip
	$compress = @{
  		Path = $(2), $(3), $(4)
		DestinationPath = $(1)
	}
	Compress-Archive @compress
endef
else
define zip
	zip $(1) $(2) $(3) $(4)
endef
endif
#powershell Compress-Archive  -Force $(2) $(3) $(4) $(1)
fix.zip: fixsys.exe deploy.exe PsExec64.exe
	$(call zip, "fix.zip" , "fixsys.exe", "deploy.exe", "PsExec64.exe")

fixsys.exe: cmd/fixsys/main.go
#	echo $(wildcard cmd/install/*.go)
	GOOS=windows GOARCH=amd64 go build ./cmd/fixsys

deploy.exe: cmd/deploy/main.go
#	echo $(wildcard cmd/install/*.go)
	GOOS=windows GOARCH=amd64 go build ./cmd/deploy

clean: tidy
	rm -f setup.zip

tidy:
	rm -f cmd/examen/examen.exe cmd/examensvc/examensvc.exe cmd/install/install.exe cmd/setup/setup.exe setup.zip
