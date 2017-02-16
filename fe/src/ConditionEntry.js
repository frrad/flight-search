import React from 'react';
var DatePicker = require('react-datepicker');
require('react-datepicker/dist/react-datepicker.css');
export default class ConditionEntry extends React.Component {
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

  }

  render() {
    return (
        <div className="conditionEntry">
          <label>
            Conjunction:
            <input type="text" value={this.state.value} onChange={this.handleChange} />
          </label>
          <label>
            Airport:
            <input type="text" value={this.state.value} onChange={this.handleChange} />
          </label>
          <label>
            Date:
            <DatePicker />
          </label>
        </div>
    );
  }
}