// Simple test to verify Jest setup is working
describe('Basic Jest functionality', () => {
  it('should pass a basic test', () => {
    expect(1 + 1).toBe(2)
  })

  it('should handle string operations', () => {
    const testString = 'CloudGate'
    expect(testString.toLowerCase()).toBe('cloudgate')
  })
}) 