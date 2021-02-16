import Router from './Router'
import Header from './common/Header.js'

import fakeSetData from './fake_set_data.json'

const App = () => {
  return (
    <div>
      <Header sets={fakeSetData}/>
      <Router />
    </div>
  );
}

export default App;
