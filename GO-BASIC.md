# golang code examples

Code examples that you can use to play with Golang

## Functions

1. **Function:**
    ```go
    func goLearn(k8s, golang string, years int) (bool, error){
        if years == 0 {
            return false, fmt.Errorf("can't lear k8s or golanf for 0 years...")
        }
        return true, nil
    ```

2. **Using the blank identifier:**
    ```go
    r, _ := goLearn("k8s", "golang", 0) //ignores the second return value
    ```

3. **Anonymous functions**
    ```go
    goLearn := func(k8s, golang string, years int) bool {
        return true
    }

    fmt.Println(goLearn("k8s", "golang", 1))

4. **Immediate invocation of a function**
    ```go
    func main() {
      func(s string) {
        fmt.Printf("I like %s", s)
      }("k8s")   // prints "I like k8s"
    }
    ```
5. **The “defer” keyword**
    ```go
    func main() {
      defer fmt.Println("k8s operator.")   // runs at the end
      fmt.Println("I can write my own")     // runs first
    }
    ```

6. **Exported / Unexported Functions**
    ```go
    func goLearn(k8s string) bool {
      return true
    }

    func GoLearn(k8s string) bool {
      return true
    }
    ```

7. **Advanced Level: Returning Functions**
    ```go
    func createAction(action string) func(string) {
      return func(article string) {
          fmt.Printf("%s, %s!\n", action, article)
      }
    }

    func main() {
       like := createAction("I like")
       hate := createAction("I hate")
        
       like("k8s")
       hate("Java")
    }
    ```

## Pointers

1. **Pointers Basics**
    ```go
    var x int = 1
    var ptr *int = &x

    fmt.Println(x)  // output: 1
    fmt.Println(ptr) // output: 0xc0000140a8
    fmt.Println(*ptr) // output: 1
    ```

2. **Memory Management with Pointers**
    ```go
    var ptr *int = new(int) // new function to allocate memory

    fmt.Println(ptr) // output: 0xc0000160c0
    fmt.Println(*ptr) // output: 0

    *ptr = 10
    fmt.Println(*ptr) // output: 10

    ptr = nil
    ```

## Structs in Golang

1. **Struct**
    ```go
    type Kubernetes struct {
        Name       string     `json:"name"`
        Version    string     `json:"version"`
        Users      []string   `json:"users,omitempty"`
        NodeNumber func() int `json:"-"`
    }

    func (k8s Kubernetes) GetUsers() {
        for _, user := range k8s.Users {
            fmt.Println(user)
        }
    }

    func (k8s *Kubernetes) AddNewUser(user string) {
        k8s.Users = append(k8s.Users, user)
    }
    ```

## Goroutines

1. **Create a Goroutine**
    ```go
    func getDeployments(name string) {
        fmt.Println(name)
    }

    func main() {
      go getDeployments("k8s")
    }
    ```

2. **Anonymous functions**
    ```go
    func main() {
       go func(name string){
         fmt.Println(name)
       }("k8s")
    }
    ```

3. **WaitGroup**
    ```go
      func main() {
        var wg sync.WaitGroup

        for i := 0; i < 5; i++ {
            wg.Add(1)
            go func(name string) {
                defer wg.Done()
            }("k8s")
        }
        wg.Wait()
      }
    ```

5. **Channels**
    ```go
    func getDeployments(name string, done chan bool) {
       done <- true
    }

    func main() {
       done := make(chan bool)

       go getDeployments("k8s", done)
       <- done
    }
    ```

6. **Channel Directions**
    ```go
    func writer(channel chan<- string, msg string) {
      channel <- msg
    }

    func reader(channel <-chan string) {
      msg := <- channel
      fmt.Println(msg)
    }
    ```

7. **Closing a channel**
    ```go
    func main() {
       ch := make(chan string)

    go func() {
        defer close(ch)
        for i := 5; i > 0; i-- {
            ch <- fmt.Sprintf("we are staring in %d sec", i)
        }
    }()

    for msg := range ch {
        fmt.Println(msg) // Receive values until the channel is closed
    }
    }
    ```
