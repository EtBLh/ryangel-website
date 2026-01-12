import Header from '../components/Header';

const About = () => {
  return (
    <div className="bg-gray-50 min-h-screen">
      <Header />
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold">About RyAngel</h1>
        <p>This is the about page.</p>
      </div>
    </div>
  );
};

export default About;