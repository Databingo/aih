# mac amd64
GOOS=windows GOARCH=386   go build -o aih_x86.exe      && zip -r ./release/aih_x86_exe.zip     aih_x86.exe       && rm -rf aih_x86.exe     
GOOS=windows GOARCH=amd64 go build -o aih_amd64.exe    && zip -r ./release/aih_amd64_exe.zip   aih_amd64.exe     && rm -rf aih_amd64.exe  
GOOS=windows GOARCH=arm   go build -o aih_arm.exe      && zip -r ./release/aih_arm_exe.zip   aih_arm.exe         && rm -rf aih_arm.exe   
GOOS=windows GOARCH=arm64 go build -o aih_arm64.exe    && zip -r ./release/aih_arm64_exe.zip   aih_arm64.exe     && rm -rf aih_arm64.exe   
GOOS=linux   GOARCH=386   go build -o aih_linux_x86    && zip -r ./release/aih_linux_x86.zip   aih_linux_x86     && rm -rf aih_linux_x86   
GOOS=linux   GOARCH=amd64 go build -o aih_linux_amd64  && zip -r ./release/aih_linux_amd64.zip aih_linux_amd64   && rm -rf aih_linux_amd64
GOOS=linux   GOARCH=arm   go build -o aih_linux_arm    && zip -r ./release/aih_linux_arm.zip aih_linux_arm       && rm -rf aih_linux_arm 
GOOS=linux   GOARCH=arm64 go build -o aih_linux_arm64  && zip -r ./release/aih_linux_arm64.zip aih_linux_arm64   && rm -rf aih_linux_arm64 
#GOOS=darwin GOARCH=386   go build -o aih_mac_x86      && zip -r ./release/aih_mac_x86.zip     aih_mac_x86       && rm -rf aih_mac_x86   
GOOS=darwin  GOARCH=amd64 go build -o aih_mac_amd64    && zip -r ./release/aih_mac_amd64.zip   aih_mac_amd64     && rm -rf aih_mac_amd64  
GOOS=darwin  GOARCH=arm64 go build -o aih_mac_arm64    && zip -r ./release/aih_mac_arm64.zip   aih_mac_arm64     && rm -rf aih_mac_arm64   
