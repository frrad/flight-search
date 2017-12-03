import React from 'react';
import ConditionEntry from './ConditionEntry.js';

export default class ResultsList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {show:'hi!'};
  }

  render() {
    return (
      <div className="searchResult">
        <div>
          <div>Depart Airpot: theairline </div>
          <div>Depart Time: theairline </div>
        </div>
        <div>
          <div>Arrival Airpot: theairline </div>
          <div>Arrival Time: theairline </div>
        </div>
        <div>Total Cost </div>
      </div>
    );
  }
}