```
        filetype := http.DetectContentType(buff)

         fmt.Println(filetype)

         switch filetype {
         case "image/jpeg", "image/jpg":
                 fmt.Println(filetype)

         case "image/gif":
                 fmt.Println(filetype)

         case "image/png":
                 fmt.Println(filetype)

         case "application/pdf":       // not image, but application !
                 fmt.Println(filetype)
         default:
                 fmt.Println("unknown file type uploaded")
         }
```

### Chat modo consola hecho en GO y PHP(https://github.com/RicardoValladares/AJAX)
```
Go: ****************************
Go: *                          *
Go: *        Go> chat_         *
Go: *                          *
Go: ****************************
Go> Hola_
```

<hr>

### Comandos para: Windows
```batch
go get github.com/RicardoValladares/chat
cd %GOPATH%/bin
chat.exe
```

### Comandos para: BSD, Linux y Mac
```bash
go get github.com/RicardoValladares/chat
cd $GOPATH/bin
./chat
```
