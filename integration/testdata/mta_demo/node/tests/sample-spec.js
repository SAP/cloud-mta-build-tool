
// use test spec file name as description to allow navigation from the test results view
describe(__filename, function () {

    it("First test", function(done) {
        expect(false).toBe(true);
        done();
    });
});