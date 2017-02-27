import React from 'react';
import ConditionEntry from './ConditionEntry.js';
import SearchTree from './SearchTree.js';

export default class SearchForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {value: '', conditions: [], conditionData: {}};

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.addCondition = this.addCondition.bind(this);
    this.onFieldChange = this.onFieldChange.bind(this);
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  onFieldChange(event) {
    this.state.conditionData[event.target.name.slice(-1)][event.target.name.slice(0,-2)] = event.target.value
  }

  handleSubmit(event) {
    event.preventDefault();
    let conditionList = [];
    for (var i = this.state.conditions.length - 1; i >= 0; i--) {
      conditionList.push(this.state.conditions[i].state)
    }
    console.log(this.state.conditionData) //this should post
  }

  addCondition(event) {
    event.preventDefault();
    const key = this.state.conditions.length.toString();
    var condition = <ConditionEntry key={key} num={key} onFieldChange={this.onFieldChange}/>;
    this.state.conditions.push(condition)
    this.state.conditionData[key] = {'conjunc':null, 'airport':null, 'date':null, 'parent':null}
    this.forceUpdate()
    console.log(this.state)
  }

  render() {
    return (
      <div className="searchForm">
        <form onSubmit={this.handleSubmit}>
          {this.state.conditions}
          <button type="button" onClick={this.addCondition}>Add Condition</button>
          <input type="submit" value="Submit" />
        </form>

        <SearchTree treeData={this.state.conditions} />
      </div>
    );
  }
}