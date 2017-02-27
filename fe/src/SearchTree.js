import React from 'react';
import ConditionEntry from './ConditionEntry.js';

export default class SearchTree extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  render() {
    let nodes = []
    for (var i = this.props.treeData.length - 1; i >= 0; i--) {
      nodes.push(<div className="treeNode" key={i} style={{paddingLeft : '20px'}}>Node</div>)
    }
    return (
      <div className="searchTree">
        {nodes}
      </div>
    );
  }
}