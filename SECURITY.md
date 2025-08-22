# Security Policy

## Supported Versions

We actively maintain security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1.0 | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow these steps:

### 1. **DO NOT** create a public GitHub issue
Security vulnerabilities should be reported privately to prevent potential exploitation.

### 2. Email Security Team
Send detailed information to: **security@wohnfair.org**

### 3. Include in Your Report
- **Description**: Clear description of the vulnerability
- **Impact**: Potential impact on users and system security
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Proof of Concept**: If available, include proof of concept code
- **Affected Versions**: Which versions are affected
- **Suggested Fix**: If you have suggestions for fixing the issue

### 4. Response Timeline
- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Depends on complexity, typically 30-90 days

### 5. Disclosure Process
- We will acknowledge receipt of your report
- We will investigate and provide regular updates
- Once fixed, we will credit you in our security advisories
- We will coordinate disclosure with you before public announcement

## Security Features

### Current Security Measures
- **Authentication**: OIDC via Keycloak with WebAuthn support
- **Authorization**: Role-based access control (RBAC)
- **Data Protection**: Encryption at rest and in transit
- **Audit Logging**: Comprehensive audit trails with cryptographic integrity
- **Input Validation**: Strict input validation and sanitization
- **Dependency Scanning**: Regular security updates for dependencies

### Planned Security Enhancements
- [ ] WebAuthn integration for device binding
- [ ] Hardware Security Module (HSM) integration
- [ ] Advanced threat detection and monitoring
- [ ] Penetration testing and security audits
- [ ] Compliance certifications (SOC 2, ISO 27001)

## Security Best Practices

### For Developers
- Follow secure coding guidelines
- Use dependency scanning tools
- Implement proper input validation
- Use parameterized queries to prevent SQL injection
- Implement proper error handling without information disclosure
- Regular security training and updates

### For Users
- Keep your authentication credentials secure
- Use strong, unique passwords
- Enable two-factor authentication when available
- Report suspicious activity immediately
- Keep your client applications updated

## Security Contacts

### Primary Security Team
- **Email**: security@wohnfair.org
- **PGP Key**: [security-pgp.asc](https://wohnfair.org/security-pgp.asc)

### Emergency Contacts
- **Critical Issues**: +49-XXX-XXXXXXX (24/7)
- **Escalation**: security-escalation@wohnfair.org

## Security Advisories

Security advisories are published at:
- [GitHub Security Advisories](https://github.com/wohnfair/wohnfair/security/advisories)
- [Security Blog](https://wohnfair.org/security)
- [Mailing List](https://groups.google.com/g/wohnfair-security)

## Responsible Disclosure

We believe in responsible disclosure and will:
- Work with researchers to understand and fix issues
- Provide appropriate credit for responsible disclosure
- Not take legal action against researchers following these guidelines
- Maintain confidentiality until issues are resolved

## Bug Bounty Program

We currently do not have a formal bug bounty program, but we do provide:
- Recognition in our security hall of fame
- Swag and merchandise for significant findings
- Potential future bug bounty program participation

## Compliance and Standards

Our security practices align with:
- **OWASP Top 10**: Web application security risks
- **NIST Cybersecurity Framework**: Risk management approach
- **GDPR**: Data protection and privacy requirements
- **EU AI Act**: AI system security and robustness
- **ISO 27001**: Information security management

## Security Updates

### Regular Updates
- **Dependencies**: Weekly security updates
- **Security Patches**: As needed, typically within 30 days
- **Major Security Releases**: Coordinated with feature releases

### Update Notifications
- **Security Advisories**: Immediate notification for critical issues
- **Release Notes**: Security updates included in release notes
- **Mailing List**: Security-focused updates and announcements

---

**Thank you for helping keep WohnFair secure!** 🔒

For general questions about WohnFair security, please use our [GitHub Discussions](https://github.com/wohnfair/wohnfair/discussions).
