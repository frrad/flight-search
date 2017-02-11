import React, { Component } from 'react';

export default class SearchForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {value: ''};

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  handleSubmit(event) {
    console.log('A thing was submitted: ' + this.state.value);
    event.preventDefault();
  }

  addCondition(event) {
    console.log('pushed the button: ' + this.state.value);
    event.preventDefault();
    var condition = <button />;

  }

  render() {
    return (
        <div className="conditionEntry">
          <label>
            Condition:
            <input type="text" value={this.state.value} onChange={this.handleChange} />
          </label>
        </div>
    );
  }
}