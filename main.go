package main

import (
  "bufio"
  "fmt"
  "os"
  "io"
  "sync"
)

var values map[string]string
var keys []string

func copyFile(src, dst string) error {
  sourceFile, err := os.Open(src)
  if err != nil {
    return err
  }
  defer sourceFile.Close()

  destinationFile, err := os.Create(dst)
  if err != nil {
    return err
  }
  defer destinationFile.Close()

  _, err = io.Copy(destinationFile, sourceFile)
  if err != nil {
    return err
  }

  return destinationFile.Sync()
}

func moveFile(src, dst string) error {
  if err := copyFile(src, dst); err != nil {
    return err
  }
  return os.Remove(src)
}

func file_changer(path string, filename string, key string, value string) {

  // Open the input file for reading
  file, err := os.Open(path)
  if err != nil {
    fmt.Println("ERROR: could not modify", path, "file : " , err)
    return
  }
  defer file.Close()

  // Create a temporary file for writing the modified content
  tempFile, err := os.CreateTemp("", filename)
  if err != nil {
    fmt.Println("ERROR: could not modify", path, "file : " , err)
    return
  }
  defer tempFile.Close()

  // Read and modify the content
  //reader := bufio.NewReader(file)
  writer := bufio.NewWriter(tempFile)
  keyLen := len(key)
  buffer := make([]byte, keyLen)

  // fmt.Println("starting to iterate over the file")

  //var buffer_str string = ""

  for {
    //n, err := reader.Read(buffer)
    n , err := file.Read(buffer)
    if n > 0 {
      // fmt.Println(string(buffer[:n]))
      if string(buffer[:n]) == (key) {
        // fmt.Println("found")
        //buffer_str += value
        writer.WriteString((value))
      } else {

        // Write only the first character of the buffer





        // Seek back in the file by keyLen - 1 bytes
        if n >= keyLen {
          // Seek back by keyLen - 1 bytes
          writer.Write(buffer[:1])
          // buffer_str += string(buffer[:1])
          //pos , _ := file.Seek(0 , io.SeekCurrent)
          //fmt.Println(pos)
          file.Seek(int64(-(keyLen - 1)), io.SeekCurrent)
          //fmt.Println(pos)
          //reader.Reset(file)
        } else {
          // If less than keyLen bytes were read, seek back by n - 1 bytes
          writer.Write(buffer[:n]) 
          //buffer_str += string(buffer[:n])
        }
        // Reset the reader to reflect the new position
        // reader.Reset(file)

      }
    }
    if err != nil {
      if err == io.EOF {
        break
      }
      fmt.Println("ERROR: reading the file")
      return
    }
  }

  // fmt.Println("starting modification on the file")

  writer.Flush()
  // reader.Flush()
  file.Close()
  tempFile.Close()
  err = os.Remove(path)
  if err != nil {
    fmt.Println("ERROR: could not modify", path, "file : " , err)
    return
  }
  // Close the temporary file to ensure all data is written

  // fmt.Println("replacing the file")

  // Replace the original file with the modified temporary file
  if err := moveFile(tempFile.Name(), path); err != nil {
    fmt.Println("ERROR: could not modify", path, "file : " , err)
    return
  }

  // fmt.Println(buffer_str)
}


func checker(wg *sync.WaitGroup , path string , filename string){
  defer wg.Done()
  for _,key := range keys{
    val := values[key]
    file_changer(path , filename , key , val);
  }
  fmt.Println("PASS : file " , path , " modified")
}


func dir_walker(path *string){
  var wg sync.WaitGroup

  files, err := os.ReadDir(*path);
  if err != nil {
    fmt.Println(err)
    return
  }

  for _, file := range files {
    real_path := *path + "/" + file.Name()
    _, err := os.ReadDir(real_path)
    if err != nil {
      // fmt.Println(real_path);
      wg.Add(1)
      go checker(&wg , real_path , file.Name())
      continue
    }

    dir_walker(&real_path);
  }

  wg.Wait()
}

func main() {
  values = make(map[string]string)
  var dir string
  fmt.Println("insert dir path : ");
  fmt.Scanln(&dir);

  for {
    fmt.Println("insert placeholder : <<END>> to break")
    var key string
    fmt.Scanln(&key)
    if key == "<<END>>" {
      break;
    }
    keys = append(keys , key)
    fmt.Println("insert value : ")
    var value string
    fmt.Scanln(&value)
    values[key] = value
  }



  dir_walker(&dir);
}
