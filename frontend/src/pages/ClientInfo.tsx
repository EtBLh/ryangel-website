import Header from '../components/Header';

const ClientInfo = () => {
  return (
    <div className="bg-gray-50 min-h-screen">
      <Header />
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold">Client Info</h1>
        <p>This is the client info page.</p>
      </div>
    </div>
  );
};

export default ClientInfo;