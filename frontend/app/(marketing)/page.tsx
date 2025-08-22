import Link from 'next/link';
import { Home, Shield, Users, BarChart3, Zap, Globe, ArrowRight, CheckCircle } from 'lucide-react';

export default function MarketingPage() {
  return (
    <div className="min-h-screen">
      {/* Navigation */}
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Home className="h-8 w-8 text-primary-600" />
              <span className="ml-2 text-xl font-bold text-gray-900">WohnFair</span>
            </div>
            <div className="hidden md:flex items-center space-x-8">
              <Link href="#features" className="text-gray-600 hover:text-gray-900">Features</Link>
              <Link href="#about" className="text-gray-600 hover:text-gray-900">About</Link>
              <Link href="#contact" className="text-gray-600 hover:text-gray-900">Contact</Link>
              <Link href="/dashboard" className="btn btn-primary btn-sm">Dashboard</Link>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="bg-gradient-to-br from-primary-50 to-secondary-50 py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            Fair Housing Allocation for{' '}
            <span className="gradient-text">Germany</span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            WohnFair implements α-fair scheduling algorithms with zero-knowledge proofs 
            to ensure equitable, transparent, and cryptographically verifiable housing distribution.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link href="/dashboard" className="btn btn-primary btn-lg">
              Get Started
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
            <Link href="#demo" className="btn btn-outline btn-lg">
              Watch Demo
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              Revolutionary Fairness Technology
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Our system combines advanced algorithms with cryptographic guarantees 
              to ensure housing allocation is both fair and transparent.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <div className="card p-6 text-center">
              <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Shield className="h-8 w-8 text-primary-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">α-Fair Scheduling</h3>
              <p className="text-gray-600">
                Advanced algorithms ensure proportional fairness while maintaining system efficiency.
              </p>
            </div>

            <div className="card p-6 text-center">
              <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Zap className="h-8 w-8 text-secondary-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Zero-Knowledge Proofs</h3>
              <p className="text-gray-600">
                Cryptographic verification without revealing sensitive personal information.
              </p>
            </div>

            <div className="card p-6 text-center">
              <div className="w-16 h-16 bg-success-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <BarChart3 className="h-8 w-8 text-success-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Real-time Analytics</h3>
              <p className="text-gray-600">
                Comprehensive metrics and fairness analysis for continuous improvement.
              </p>
            </div>

            <div className="card p-6 text-center">
              <div className="w-16 h-16 bg-warning-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Users className="h-8 w-8 text-warning-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Group Equity</h3>
              <p className="text-gray-600">
                Special consideration for refugees, disabled, seniors, and low-income groups.
              </p>
            </div>

            <div className="card p-6 text-center">
              <div className="w-16 h-16 bg-error-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Globe className="h-8 w-8 text-error-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Multi-City Support</h3>
              <p className="text-gray-600">
                Scalable architecture supporting Berlin, Munich, and other German cities.
              </p>
            </div>

            <div className="card p-6 text-center">
              <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <CheckCircle className="h-8 w-8 text-primary-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">GDPR Compliant</h3>
              <p className="text-gray-600">
                Built with privacy by design, meeting all European data protection requirements.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              How WohnFair Works
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              A simple three-step process ensures fair and transparent housing allocation.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="w-20 h-20 bg-primary-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-2xl font-bold">
                1
              </div>
              <h3 className="text-xl font-semibold mb-3">Submit Request</h3>
              <p className="text-gray-600">
                Users submit housing requests with preferences, urgency levels, and group classification.
              </p>
            </div>

            <div className="text-center">
              <div className="w-20 h-20 bg-secondary-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-2xl font-bold">
                2
              </div>
              <h3 className="text-xl font-semibold mb-3">Fair Processing</h3>
              <p className="text-gray-600">
                Our α-fair algorithm processes requests considering group weights, urgency, and wait times.
              </p>
            </div>

            <div className="text-center">
              <div className="w-20 h-20 bg-success-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-2xl font-bold">
                3
              </div>
              <h3 className="text-xl font-semibold mb-3">Verifiable Allocation</h3>
              <p className="text-gray-600">
                Results are cryptographically verified and can be audited without compromising privacy.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* About Section */}
      <section id="about" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="text-4xl font-bold text-gray-900 mb-6">
                Built for German Municipalities
              </h2>
              <p className="text-lg text-gray-600 mb-6">
                WohnFair is designed specifically for German housing authorities, 
                addressing the unique challenges of urban housing allocation while 
                maintaining the highest standards of fairness and transparency.
              </p>
              <div className="space-y-4">
                <div className="flex items-center">
                  <CheckCircle className="h-5 w-5 text-success-600 mr-3" />
                  <span>Compliant with German housing laws and regulations</span>
                </div>
                <div className="flex items-center">
                  <CheckCircle className="h-5 w-5 text-success-600 mr-3" />
                  <span>Multi-language support (German, English, Turkish)</span>
                </div>
                <div className="flex items-center">
                  <CheckCircle className="h-5 w-5 text-success-600 mr-3" />
                  <span>Integration with existing municipal systems</span>
                </div>
                <div className="flex items-center">
                  <CheckCircle className="h-5 w-5 text-success-600 mr-3" />
                  <span>Local data sovereignty and privacy protection</span>
                </div>
              </div>
            </div>
            <div className="bg-gradient-to-br from-primary-100 to-secondary-100 rounded-lg p-8">
              <h3 className="text-2xl font-semibold mb-4">Technical Highlights</h3>
              <ul className="space-y-3 text-gray-700">
                <li>• α-fair scheduling with configurable fairness parameters</li>
                <li>• Zero-knowledge proofs using Halo2/PLONK</li>
                <li>• Real-time fairness metrics and monitoring</li>
                <li>• Scalable microservices architecture</li>
                <li>• Comprehensive audit logging and compliance</li>
                <li>• Machine learning for demand prediction</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-primary-600">
        <div className="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
          <h2 className="text-4xl font-bold text-white mb-6">
            Ready to Transform Housing Allocation?
          </h2>
          <p className="text-xl text-primary-100 mb-8">
            Join German municipalities in implementing fair, transparent, and 
            cryptographically verifiable housing allocation systems.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link href="/dashboard" className="btn bg-white text-primary-600 hover:bg-gray-100 btn-lg">
              Start Free Trial
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
            <Link href="#contact" className="btn btn-outline border-white text-white hover:bg-white hover:text-primary-600 btn-lg">
              Contact Sales
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center mb-4">
                <Home className="h-8 w-8 text-primary-400" />
                <span className="ml-2 text-xl font-bold">WohnFair</span>
              </div>
              <p className="text-gray-400">
                Building fair housing allocation for the digital age.
              </p>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Product</h4>
              <ul className="space-y-2 text-gray-400">
                <li><Link href="#features" className="hover:text-white">Features</Link></li>
                <li><Link href="/dashboard" className="hover:text-white">Dashboard</Link></li>
                <li><Link href="#pricing" className="hover:text-white">Pricing</Link></li>
                <li><Link href="#api" className="hover:text-white">API</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Company</h4>
              <ul className="space-y-2 text-gray-400">
                <li><Link href="#about" className="hover:text-white">About</Link></li>
                <li><Link href="#careers" className="hover:text-white">Careers</Link></li>
                <li><Link href="#blog" className="hover:text-white">Blog</Link></li>
                <li><Link href="#press" className="hover:text-white">Press</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Support</h4>
              <ul className="space-y-2 text-gray-400">
                <li><Link href="#help" className="hover:text-white">Help Center</Link></li>
                <li><Link href="#contact" className="hover:text-white">Contact</Link></li>
                <li><Link href="#status" className="hover:text-white">System Status</Link></li>
                <li><Link href="#docs" className="hover:text-white">Documentation</Link></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <p>&copy; 2024 WohnFair. All rights reserved. | Privacy Policy | Terms of Service</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
