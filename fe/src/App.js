import React, { Component } from 'react';
import SearchForm from './searchform.js';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <h2>Flight Search</h2>
          <h3></h3>
        </div>
        <SearchForm />
      </div>
    );
  }
}

export default App
