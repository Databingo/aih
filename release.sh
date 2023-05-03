# mac amd64
GOOS=windows GOARCH=386   go build -o aih.exe   && zip -r ./release/aih_x86_exe.zip     aih.exe     && rm -rf aih.exe   # aih_x86.exe    
GOOS=windows GOARCH=amd64 go build -o aih.exe   && zip -r ./release/aih_amd64_exe.zip   aih.exe     && rm -rf aih.exe   # aih_amd64.exe  
GOOS=windows GOARCH=arm   go build -o aih.exe   && zip -r ./release/aih_arm_exe.zip     aih.exe     && rm -rf aih.exe   # aih_arm.exe   
GOOS=windows GOARCH=arm64 go build -o aih.exe   && zip -r ./release/aih_arm64_exe.zip   aih.exe     && rm -rf aih.exe   # aih_arm64.exe  
GOOS=linux   GOARCH=386   go build -o aih       && zip -r ./release/aih_linux_x86.zip   aih         && rm -rf aih       # aih_linux_x86  
GOOS=linux   GOARCH=amd64 go build -o aih       && zip -r ./release/aih_linux_amd64.zip aih         && rm -rf aih       # aih_linux_amd64
GOOS=linux   GOARCH=arm   go build -o aih       && zip -r ./release/aih_linux_arm.zip   aih         && rm -rf aih       # aih_linux_arm 
GOOS=linux   GOARCH=arm64 go build -o aih       && zip -r ./release/aih_linux_arm64.zip aih         && rm -rf aih       # aih_linux_arm64
#GOOS=darwin GOARCH=386   go build -o aih       && zip -r ./release/aih_mac_x86.zip     aih         && rm -rf aih       # aih_mac_x86   
GOOS=darwin  GOARCH=amd64 go build -o aih       && zip -r ./release/aih_mac_amd64.zip   aih         && rm -rf aih       # aih_mac_amd64  
GOOS=darwin  GOARCH=arm64 go build -o aih       && zip -r ./release/aih_mac_arm64.zip   aih         && rm -rf aih       # aih_mac_arm64  
