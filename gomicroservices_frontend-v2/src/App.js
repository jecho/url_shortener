import React, { Component } from 'react';

class URLForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {result: 'Nenley stomach hurts... Feed it URLs please!', value: 'Enter Input Here'};
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  handleSubmit(event) {
    event.preventDefault();
    alert('submitted: ' + this.state.value);

    fetch('/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        'url' : this.state.value
      })
    }).then((response) => response.json())
    .then((responseJson) => {
      this.setState({result: JSON.stringify(responseJson)})
    })
    .catch((error) => {
      console.error(error);
    });
  }

  render() {
    return (
      <div>
        <form id="formoid" method="POST" onSubmit={this.handleSubmit}>
          <label>
            <input type="text" value={this.state.value} onChange={this.handleChange} />
          </label>
          <input type="submit" value="Submit" />
        </form>
        <br/>
        <p className="set-in">{this.state.result}</p>
      </div>
    );
  }
}


class App extends Component {
  render() {
    return (
      <div id="main">
        <style>{'body { background-color: #c87548; }'}</style>
        <div className="cubano">
          <div className="wrap free">
            <h1 className="title">That's So Meta!</h1>
              <p>The intellect of the wise is like glass; it admits the link through the glass and reflects on its shortening. But its similar than that, feed Nenley.</p>
              <p><q>Alright, now that you got it, let the orange cookie monster Nenley, eat it!</q></p>
              <p>signed. <b><i>`foxley</i></b></p>
              <URLForm />
          </div>
        </div>
      </div>
    );
  }
}

export default App;
