import React from 'react';
import DatePicker from 'react-datepicker'
require('react-datepicker/dist/react-datepicker.css');
export default class ConditionField extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <input name={this.props.name} type="text" value={this.props.value} conditionNumber={this.props.conditionNumber} onChange={this.props.onChange} />
    );
  }
}