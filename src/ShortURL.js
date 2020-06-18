import React from 'react';

class ShortURL extends React.Component{

    render(){
      console.log("ShortURL_render() called",this.props)
        if (this.props.activeState === "Valid URL" && this.props.clicked === "True" && this.props.key1 != null){
            //createNExecGetReq(document.getElementById("formURL").value,this)
            var key = this.props.key1
            return(
            //<p>Valid URL {this.state.key}</p>,
            <p>NEW Valid URL {this.props.myURL.slice(0,this.props.myURL.length - 1) + key}</p>,
            <a href={this.props.myURL.slice(0,this.props.myURL.length - 1) + key}>{this.props.myURL.slice(0,this.props.myURL.length - 1) + key}</a>
            )
        } else if (this.props.activeState === "Invalid URL" && this.props.clicked === "True"){
            return (
            <p>Invalid URL</p>
            )
        } else {
          return (
            null
          )
        }
    }


}

export default ShortURL