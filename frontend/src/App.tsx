import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/client/Home';
import About from './pages/client/About';
import Checkout from './pages/client/Checkout';
import ClientInfo from './pages/client/ClientInfo';
import GoogleCallback from './pages/client/GoogleCallback';
import ProductInfo from './pages/client/ProductInfo';

import './index.css'
import ClientLayout from './pages/client/ClientLayout';
import { Toaster } from 'sonner';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<ClientLayout />}>
          <Route index element={<Home />} />
          <Route path="about" element={<About />} />
          <Route path="checkout" element={<Checkout />} />
          <Route path="client-info" element={<ClientInfo />} />
          <Route path="google-callback" element={<GoogleCallback />} />
          <Route path="/product/:productId" element={<ProductInfo />} />
        </Route>
      </Routes>
      <Toaster />
    </BrowserRouter>
  );
}

export default App;
