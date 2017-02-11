import React, { Component } from 'react';
import ConditionEntry from './searchform.js';

export default class SearchForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {value: '', conditions: ['one']};

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
    let conditions = []
    this.state.conditions.forEach(function(con) {
        console.log(con)
        conditions.push(<ConditionEntry />)

    })
    return (
      <form onSubmit={this.handleSubmit}>
        <div className="conditionEntry">
          {conditions}
          <label>
            Condition:
            <input type="text" value={this.state.value} onChange={this.handleChange} />
          </label>
        </div>
        <button onClick="{addCondition()}">
          Add Condition
        </button>
        <input type="submit" value="Submit" />
      </form>
    );
  }
}