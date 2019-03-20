package main

import (
  "net/http"
  "net/url"
  "fmt"
  "encoding/json"
  "io/ioutil"
)

func check(e error) {
  if e!= nil {
    panic(e)
  }
}

func getRequest(w http.ResponseWriter, r *http.Request) {
  byteBody, err := ioutil.ReadAll(r.Body)
  check(err)
  var data map[string]interface{}
  
  err = json.Unmarshal(byteBody, &data)
  check(err)

  switch data["object_kind"] {
    case "push":
      generatePushMsg(data)
    case "tag_push":
      generatePushTagMsg(data)
    case "issue":
      generateIssueMsg(data)
    case "note":
      generateNoteMsg(data)
    case "merge_request":
      generateMergeMsg(data)
    case "wiki_page":
      generateWikiMsg(data)
    case "pipeline":
      generatePipelineMsg(data)
    case "build":
      generateBuildMsg(data)
    default:
      fmt.Println("No generazione :(")
  }
}

func generatePushMsg(data map[string]interface{}) {
  var s []interface{}
  commits := data["commits"].([]interface{})
  repository := data["repository"].(map[string]interface{})
  for i := 0; i < len(commits); i++ {
    currentCommit := commits[i].(map[string]interface{})
    s = append(s, currentCommit["message"])
  }
  
  u := fmt.Sprintf("%v", data["user_name"])
  r := fmt.Sprintf("%v", repository["homepage"])
  c := fmt.Sprintf("%v", s)
  msg := "New push to GitLab!\n" + "User: " + u + "\n"  + "Commit(s): " + c + "\n" + "Repository: " + r
  sendToBot(msg)

}

func generateIssueMsg(data map[string]interface{}) {
  repository := data["repository"].(map[string]interface{})
  attributes := data["object_attributes"].(map[string]interface{})
  user := data["user"].(map[string]interface{})

  u := fmt.Sprintf("%v", user["name"])
  it := fmt.Sprintf("%v", attributes["title"])
  id := fmt.Sprintf("%v", attributes["description"])
  r := fmt.Sprintf("%v", repository["homepage"])
  msg := "An issue was created/edited/closed!\n" + "User: " + u + "\n" + "Issue title: " + it + "\n" + "Issue description: " + id + "\n" + "Respository: " + r
  sendToBot(msg)

}

func generateMergeMsg(data map[string]interface{}) {
  attributes := data["object_attributes"].(map[string]interface{})
user := data["user"].(map[string]interface{})
  target := attributes["target"].(map[string]interface{})
  
  u := fmt.Sprintf("%v", user["name"])
  mt := fmt.Sprintf("%v", attributes["title"])
  ms := fmt.Sprintf("%v", attributes["state"])
  r := fmt.Sprintf("%v", target["web_url"])
  msg := "A merge request was created/updated/merged/closed!\n" + "User: " + u + "\n" + "MR title: " + mt + "\n" + "MR state: " + ms + "\n" + "Repository: " + r 
  sendToBot(msg)
}

func generatePipelineMsg(data map[string]interface{}) {
  attributes := data["object_attributes"].(map[string]interface{})
  user := data["user"].(map[string]interface{})
  project := data["project"].(map[string]interface{})

  u := fmt.Sprintf("%v", user["name"])
  s := fmt.Sprintf("%v", attributes["status"])
  r := fmt.Sprintf("%v", project["web_url"])
  msg := "The status of a pipeline has changed!" + "\n" + "User: " + u + "\n" + "Status: " +  s + "\n" + "Repository: " + r
  sendToBot(msg)

}

func generateBuildMsg(data map[string]interface{}) {
  repository := data["repository"].(map[string]interface{})
  user := data["user"].(map[string]interface{})

  u := fmt.Sprintf("%v", user["name"])
  bn := fmt.Sprintf("%v", data["build_name"])
  bs := fmt.Sprintf("%v", data["build_status"])
  r := fmt.Sprintf("%v", repository["homepage"])
  msg := "The status of a build has changed!" + "\n" + "User: " + u + "\n" + "Build name: " + bn + "\n"  + "Build status: " +  bs + "\n" + "Repository: " + r
  sendToBot(msg)

}

func generateWikiMsg(data map[string]interface{}) {
  project := data["project"].(map[string]interface{})
  user := data["user"].(map[string]interface{})
  attributes := data["object_attributes"].(map[string]interface{})
  wiki := data["wiki"].(map[string]interface{})

  u := fmt.Sprintf("%v", user["name"])
  w := fmt.Sprintf("%v", wiki["web_url"]  )
  t := fmt.Sprintf("%v", attributes["title"])
  r := fmt.Sprintf("%v", project["web_url"])
  msg := "A wiki page was created/updated/deleted!" + "\n" + "User: " + u + "\n" + "Wiki: " + w + "\n"  + "Page title: " +  t + "\n" + "Repository: " + r
  sendToBot(msg)

}

func generatePushTagMsg(data map[string]interface{}) {
  repository := data["repository"].(map[string]interface{})
  
  u := fmt.Sprintf("%v", data["user_name"])
  t := fmt.Sprintf("%v", data["ref"])
  r := fmt.Sprintf("%v", repository["homepage"])
  msg := "A tag was created/deleted!" + "\n" + "User: " + u + "\n" + "Tag: " + t + "\n" + "Repository: " + r
  sendToBot(msg)

}

// TODO Print which commit was annotated
func generateNoteMsg(data map[string]interface{}) {
  attributes := data["object_attributes"].(map[string]interface{})
  repository := data["repository"].(map[string]interface{})
  user := data["user"].(map[string]interface{})

  u := fmt.Sprintf("%v", user["name"])
  n := fmt.Sprintf("%v", attributes["note"])
  r := fmt.Sprintf("%v", repository["homepage"])
  
  switch attributes["noteable_type"] {
    case "Commit":
      msg := "A commit was annotated!" + "\n" + "User: " + u + "\n" + "Note: " + n + "\n" + "Repository: " + r
      sendToBot(msg)
    case "MergeRequest":
      msg := "A merge request was annotated!" + "\n" + "User: " + u + "\n" + "Note: " + n + "\n" + "Repository: " + r
      sendToBot(msg)
    case "Issue":
      msg := "An issue was annotated!" + "\n" + "User: " + u + "\n" + "Note: " + n + "\n" + "Repository: " + r
      sendToBot(msg)
    case "Snippet":
      msg := "A code snippet was annotated!" + "\n" + "User: " + u + "\n" + "Note: " + n + "\n" + "Repository: " + r
      sendToBot(msg)
    default:
      fmt.Println("Something went wrong :(")
    }
}

func sendToBot(msg string) {
	var link string = ""

  chat_id, err := ioutil.ReadFile("config/chat_id")
  check(err)
  token, err := ioutil.ReadFile("config/token")
  check(err)

  t:= &url.URL{Path: msg}
  encodedMsg := t.String()
  
  if encodedMsg[:2] == "./" {
    link = "https://api.telegram.org/bot" + string(token)[:len(token)-1] + "/sendMessage?chat_id=" + string(chat_id)[:len(chat_id)-1]  + "&text=" + encodedMsg[2:]
  } else {
    link = "https://api.telegram.org/bot" + string(token)[:len(token)-1] + "/sendMessage?chat_id=" + string(chat_id)[:len(chat_id)-1]  + "&text=" + encodedMsg
  }
  
  fmt.Println(link)
  
  resp, err := http.Get(link)
  check(err)
  fmt.Println(resp)
}

func main() {
  port, err := ioutil.ReadFile("config/port")
  check(err)
  http.HandleFunc("/", getRequest)
  if err := http.ListenAndServe(":" + string(port)[:len(port)-1], nil); err != nil {
    panic(err)
  }
}
