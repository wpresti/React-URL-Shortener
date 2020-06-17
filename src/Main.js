import React from 'react';
import logo from './logo.svg';
import './App.css';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form'
import ShortURL from './ShortURL';
class Main extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      activeState: "",
      clicked: "",
      key: null,
      myURL: window.location.href,
      redirectURL: null
    };
  }
  componentDidMount(){
    console.log("component mounted!!!")
    //redirect check -- refactor into function ASAP
    if(window.location.href !== "http://localhost:3000/"){
      console.log("not equal -- redirecting...")
      console.log(window.location.pathname)
      //call func to lookup key, strip leading / first
      //in func above^ return fullurl. hardcoded for now.

      lookupURLFromKey(window.location.pathname,this)
      if (this.state.redirectURL !== null){
        console.log("i should be redirecrting now...")
        window.location.href = this.state.redirectURL
      }
      //window.location.href = "http://google.com";
    }
  }
    //this.handleClick.bind(this);
    //this.handleClick = this.handleClick.bind(this);
  handleClick = () => {
    console.log("Button clicked!")
    var URL = document.getElementById("formURL").value
    console.log(URL)
    console.log(validURL(URL))
    //this.setState({clicked: "True"})
    if (validURL(URL) === true) {
      console.log("Valid URL")
      this.setState({clicked: "True", activeState:"Valid URL"})
      createNExecPutReq(document.getElementById("formURL").value,this)


    } else{
      //not valid URL
      console.log("not valid URL")
      this.setState({clicked: "True", activeState:"Invalid URL"})
    }
  }
  

  render(){
    return (
      <div className="App">
        
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/Main.js</code> and save to reload.
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>

          <Form>
            <Form.Group controlId="formURL">
              <Form.Label>URL to shorten</Form.Label>
              <Form.Control type="url" placeholder="Enter URL" />
            </Form.Group>
          </Form>

          <Button variant="primary" size="lg" onClick={this.handleClick}>Primary</Button>{' '}
          {/* <p> {this.state.activeState}</p> */}
        

        <ShortURL activeState={this.state.activeState} clicked={this.state.clicked} key1={this.state.key} myURL={this.state.myURL}/>
        </header>
      </div>
    )
  }
}


function validURL(str) {
  var pattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
    '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
    '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
    '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
    '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
    '(\\#[-a-z\\d_]*)?$','i'); // fragment locator
  return !!pattern.test(str);
}

//rename to putReq
function createNExecPutReq(url,self){
  var data = {}
  data.longURL = url
  var json = JSON.stringify(data)
  var req = new XMLHttpRequest();
  req.open("PUT","http://localhost:8080/",true)
  req.setRequestHeader('Content-Type','application/json')
  // https://stackoverflow.com/questions/44304773/ionic-2-calling-function-after-xhr-onload
  req.onload = async () => {
    var z = await JSON.parse(req.responseText)
    console.log("inside body", z)
    //set state in here
    self.setState({key: "/"+z.shortURL})
  }
  req.send(json)
  
}

function lookupURLFromKey(key,self){
  var keySplit = key.split("/");
  console.log("key", keySplit[1])
  createNExecGetReq(keySplit[1],self)
  
  
}
// post is acting as a get request with payload bc xmlhttpreq() doesn't send body on GET request
function createNExecGetReq(key,self){
  var data = {}
  data.shortURL = key
  var json = JSON.stringify(data)
  var req = new XMLHttpRequest()
  console.log("json payload:", json)
  req.open("POST","http://localhost:8080/",true)
  req.setRequestHeader('Content-Type','application/json')
  req.onload = async () => {
    console.log("response", req.response)
    var jsonBody = await JSON.parse(req.responseText)
    console.log("get json data:", jsonBody)
    //set state below
    // jsonBody.URL
    //console.log(this,self)
    self.setState({redirectURL: jsonBody.URL})
    await new Promise(r => setTimeout(r, 2000));
    window.location.href = jsonBody.URL
    
  }
  req.send(json)

}


export default Main;
