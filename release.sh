# mac amd64
GOOS=windows GOARCH=386   go build -o aih_32.exe       && zip -r ./release/aih_32_exe.zip      aih_32.exe        
GOOS=linux   GOARCH=386   go build -o aih_linux_32     && zip -r ./release/aih_linux_32.zip    aih_linux_32     
GOOS=darwin  GOARCH=386   go build -o aih_mac_32       && zip -r ./release/aih_mac_32.zip      aih_mac_32       
GOOS=windows GOARCH=amd64 go build -o aih_amd64.exe    && zip -r ./release/aih_amd64_exe.zip   aih_amd64.exe    
GOOS=linux   GOARCH=amd64 go build -o aih_linux_amd64  && zip -r ./release/aih_linux_amd64.zip aih_linux_amd64  
GOOS=darwin  GOARCH=amd64 go build -o aih_mac_amd64    && zip -r ./release/aih_mac_amd66.zip   aih_mac_amd64    
GOOS=windows GOARCH=arm64 go build -o aih_arm64.exe    && zip -r ./release/aih_arm64_exe.zip   aih_arm64.exe      
GOOS=linux   GOARCH=arm64 go build -o aih_linux_arm64  && zip -r ./release/aih_linux_arm64.zip aih_linux_arm64    
GOOS=darwin  GOARCH=arm64 go build -o aih_mac_arm64    && zip -r ./release/aih_mac_arm64.zip   aih_mac_arm64      
