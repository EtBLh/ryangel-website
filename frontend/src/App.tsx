import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import About from './pages/About';
import Checkout from './pages/Checkout';
import ClientInfo from './pages/ClientInfo';
import ProductInfo from './pages/ProductInfo';

import './index.css'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/checkout" element={<Checkout />} />
        <Route path="/client-info" element={<ClientInfo />} />
        <Route path="/product/:productId" element={<ProductInfo />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
