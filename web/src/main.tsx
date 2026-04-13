/**
 * Application Entry Point
 *
 * Initializes the React application and renders the root App component.
 * Uses StrictMode for development-time checks.
 */

import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

// Mount the React application to the DOM
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
