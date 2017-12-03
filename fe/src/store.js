import 'whatwg-fetch'
import 'mobx'

export default class StateStore {
    fetchSearch(query) {
        console.log('search');
        fetch('/search', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              query: query,
              test: 'test',
            })
        }).then((response) => console.log(response)).catch((error) => console.log(error))
    }
}