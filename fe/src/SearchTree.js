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
    for (var i = Object.keys(this.props.treeData).length - 1; i >= 0; i--) {
      let node = this.props.treeData[i];
      console.log(node)
      nodes.push(
        <li>
          <div className="treeNode" key={i} style={{paddingLeft : '20px'}}>
            <span>Node</span>
            <span className="nodeData">{node.conjunc}</span>
            <span className="nodeData">{node.date}</span>
            <span className="nodeData">{node.airport}</span>
            <span className="nodeData">{node.parentId}</span>
          </div>
        </li>)
    }
    return (
    <div>
      <div className="tree">
        <ul>
          {nodes}
        </ul>
      </div>
    </div>
    );
  }
}