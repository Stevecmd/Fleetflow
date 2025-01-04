import React, { useState } from 'react';
import Logo from './../assets/logo.svg';

/**
 * The LandingPage component renders the main landing page of the application.
 * It includes a navigation bar with links to the features, about, and contact pages.
 * It also includes a hero section, a features section, an about section, a call to action section, and a footer.
 */
const LandingPage: React.FC = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    return (
        <div className="bg-gray-50">
            {/* Navbar */}
            <nav className="bg-white shadow-lg">
                <div className="max-w-7xl mx-auto px-4">
                    <div className="flex justify-between h-16">
                        {/* Logo/Brand */}
                        <div className="flex items-center">
                            <a href="/" className="flex-shrink-0">
                                <img 
                                    src={Logo}
                                    alt="FleetFlow Logo"
                                    className="h-8 w-auto sm:h-10"
                                />
                            </a>
                        </div>
                        
                        {/* Desktop Navigation */}
                        <div className="hidden md:flex items-center space-x-8">
                            <a href="#features" className="text-gray-600 hover:text-indigo-600">Features</a>
                            <a href="#about" className="text-gray-600 hover:text-indigo-600">About</a>
                            <a href="#contact" className="text-gray-600 hover:text-indigo-600">Contact</a>
                            <a href="/login" className="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700">
                                Login
                            </a>
                        </div>

                        {/* Mobile menu button */}
                        <div className="md:hidden flex items-center">
                            <button onClick={() => setIsMenuOpen(!isMenuOpen)} className="text-gray-600 hover:text-gray-900">
                                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    {isMenuOpen ? (
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                    ) : (
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                                    )}
                                </svg>
                            </button>
                        </div>
                    </div>

                    {/* Mobile Navigation */}
                    {isMenuOpen && (
                        <div className="md:hidden pb-4">
                            <div className="flex flex-col space-y-3">
                                <a href="#features" className="text-gray-600 hover:text-indigo-600">Features</a>
                                <a href="#about" className="text-gray-600 hover:text-indigo-600">About</a>
                                <a href="#contact" className="text-gray-600 hover:text-indigo-600">Contact</a>
                                <a href="/login" className="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 text-center">
                                    Login
                                </a>
                            </div>
                        </div>
                    )}
                </div>
            </nav>
            {/* Hero Section */}
            <section className="text-center py-20 bg-hero">
                <div className="bg-black bg-opacity-45 p-6 rounded-lg inline-block">
                    <h1 className="text-4xl font-bold text-white">Revolutionize Your Fleet Management</h1>
                    <p className="mt-4 text-lg text-white">Efficient, Reliable, and User-Friendly Solutions</p>
                    <a href="#features" className="mt-6 inline-block px-6 py-3 bg-indigo-600 text-white rounded-lg shadow hover:bg-indigo-700 transition">
                        Get Started
                    </a>
                </div>
            </section>

            {/* Features Section */}
            <section id="features" className="py-20">
                <h2 className="text-3xl font-bold text-center text-gray-900">Our Features</h2>
                <div className="mt-10 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    <div className="bg-white p-6 rounded-lg shadow feature-card">
                        <h3 className="text-xl font-semibold">Driver Management</h3>
                        <p className="mt-2 text-gray-600">Manage driver profiles, track performance, and ensure compliance.</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow feature-card">
                        <h3 className="text-xl font-semibold">Vehicle Fleet Management</h3>
                        <p className="mt-2 text-gray-600">Track vehicle status, maintenance schedules, and analytics.</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow feature-card">
                        <h3 className="text-xl font-semibold">Delivery Operations</h3>
                        <p className="mt-2 text-gray-600">Optimize routes and track deliveries in real-time.</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow feature-card">
                        <h3 className="text-xl font-semibold">Warehouse Management</h3>
                        <p className="mt-2 text-gray-600">Monitor inventory and coordinate loading/unloading.</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow feature-card">
                        <h3 className="text-xl font-semibold">User Roles</h3>
                        <p className="mt-2 text-gray-600">Manage access and permissions for different user roles.</p>
                    </div>
                </div>
            </section>

            {/* About Section */}
            <section id="about" className="py-20 bg-gray-100">
                <h2 className="text-3xl font-bold text-center text-gray-900">About Us</h2>
                <p className="mt-4 text-lg text-center text-gray-600">FleetFlow is dedicated to providing top-notch fleet management solutions.</p>
                <p className="mt-2 text-center text-gray-600">Our mission is to help businesses optimize their fleet operations and improve efficiency.</p>
            </section>

            {/* Call to Action Section */}
            <section id="contact" className="text-center py-20 bg-indigo-600 text-white">
                <h2 className="text-3xl font-bold">Ready to Get Started?</h2>
                <p className="mt-4">Join us today and take control of your fleet management.</p>
                <a href="#contact" className="mt-6 inline-block px-6 py-3 bg-white text-indigo-600 rounded-lg shadow hover:bg-gray-200 transition">
                    Contact Us
                </a>
            </section>

            {/* Footer */}
            <footer className="py-10 text-center text-gray-600">
                <p>&copy; {new Date().getFullYear()} FleetFlow. All rights reserved.</p>
            </footer>
        </div>
    );
};

export default LandingPage;