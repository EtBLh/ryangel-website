import Header from '../components/Header';

const Checkout = () => {
  return (
    <div className="bg-gray-50 min-h-screen">
      <Header />
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold">Checkout</h1>
        <p>This is the checkout page.</p>
      </div>
    </div>
  );
};

export default Checkout;