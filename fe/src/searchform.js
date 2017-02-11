import React from 'react';
import ConditionEntry from './ConditionEntry.js';

export default class SearchForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {value: '', conditions: [{}]};

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.addCondition = this.addCondition.bind(this);
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
    var condition = {};
    this.state.conditions.push(condition)
    this.forceUpdate()
    console.log(this.state)

  }

  render() {
    console.log('rendiner')
    let conditionElements = []
    console.log(this.state.conditions)
    let i = 1
    this.state.conditions.forEach(function(con) {
        i = i + 1
        conditionElements.push(<ConditionEntry key={i.toString()}/>)
    })
    return (
      <form onSubmit={this.handleSubmit}>
        {conditionElements}
        <button type="button" onClick={this.addCondition}>Add Condition</button>
        <input type="submit" value="Submit" />
      </form>
    );
  }
}