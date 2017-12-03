import React from 'react';
import { inject, observer } from 'mobx-react';
import DatePicker from 'react-datepicker';
import StateStore from './store.js';
require('react-datepicker/dist/react-datepicker.css');
import ConditionField from './ConditionField.js';
export default class ConditionEntry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      conjunc: '',
      airport: '',
      date: '',
      parentId: '',
    };

    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event) {
    let name = event.target.name;
    let value =  event.target.value;
    if (name.indexOf('conjunc') >= 0) {
      this.setState({'conjunc':value})
    }
    else if (name.indexOf('airport') >= 0) {
      this.setState({'airport':value})
    }
    else if (name.indexOf('date') >= 0) {
      this.setState({'date':value})
    }
    else if (name.indexOf('parent') >= 0) {
      this.setState({'parentId':value})
    }
    this.props.onFieldChange(event);
  }

  render() {
    return (
        <div className="conditionEntry">
          <label>Condition ID: {this.props.num} </label>
          <label>
            ParentId:
            <ConditionField name={'parent-'+this.props.num} value={this.state.parent} conditionNumber={this.props.num} onChange={this.handleChange} />
          </label>
          <label>
            Node Type:
            <ConditionField name={'conjunc-'+this.props.num} value={this.state.conjunc} conditionNumber={this.props.num} onChange={this.handleChange} />
          </label>
          <label>
            Airport:
            <ConditionField name={'airport-'+this.props.num} value={this.state.airport} conditionNumber={this.props.num} onChange={this.handleChange} />
          </label>
          <label>
            Date:
            <ConditionField name={'date-'+this.props.num} value={this.state.date} conditionNumber={this.props.num} onChange={this.handleChange} />
          </label>
        </div>
    );
  }
}