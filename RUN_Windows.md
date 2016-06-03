# Installation of Elasticbeat on Windows

In order to install Elasticbeat on Windows, refer to the following steps:

1. Download Go from official Go language site using the following link:
https://storage.googleapis.com/golang/go1.6.2.windows-amd64.msi
Upon opening the link click it will download the MSI distributable package.

2. Extract the downloaded zip package by either unzipping it using WinRar, 7Zip etc.( If you don't have download it.).
It will extract the files and folder in directory. For simplicity extract them at C:\Go.

3. Now for running Go in windows, you need to set path of Go in environment variable settings of Windows. Firstly open properties of My Computer. Select Advanced system settings and then click on Advanced tab wherein you will click environment variables options.
   
   Then double click Path variable (under System variables) and move towards the end of text box insert a semi colon if not inserted and add the location of bin folder of Go such as: C:\Go\bin. Then click on ok to all the windows opened.

  **Note: Do not delete anything within the path variable textbox.**

4. To validate whether Go is successfully installed type the following command in command prompt:
   go version

   It will print the version of Go language as: go version go1.6.2 windows/amd64

5. Create another system variable named GOPATH referencing the workspace of GO. For example our workspace is present inside C:\Go

   After opening Environment variables click on New (under System variables) and give variable name as GOPATH and variable value as C:\Go.

6. Navigate to the workspace and create a directory named src/github.com/radoondas

7. Navigate to the directory and clone the elasticbeat repository or download the elasticbeat repository in the created directory and extract it.

8. Build elasticbeat by opening command prompt in the navigated directory ($GOPATH/src/github.com/radoondas/elasticbeat) by pressing Shift button and right click in the folder to select an option of Open Command Window Here. Use following command in the command prompt:-
   ```bash
   go build
   ```
   It will create a executable file named elasticbeat.exe.

   **Note: If you are facing errors after running make command, make sure that GO installed is greater than version 1.5.**

9. Run elasticbeat. by using following command in the command prompt:-
```bash
elasticbeat.exe -e -v -d elasticbeat -c elasticbeat.yml
```
This will run elasticbeat successfully and will start gathering information related to the Elasticsearch nodes and indices them in Elasticsearch as per it’s configuration file.
