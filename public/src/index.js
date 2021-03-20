import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './App';
import { BrowserRouter } from 'react-router-dom'
import reportWebVitals from './reportWebVitals';
import SimpleReactLightbox from 'simple-react-lightbox'

import 'react-bulma-components/dist/react-bulma-components.min.css';

ReactDOM.render(
  <React.StrictMode>
    <SimpleReactLightbox>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </SimpleReactLightbox>
  </React.StrictMode>,
  document.getElementById('root')
);

reportWebVitals();
