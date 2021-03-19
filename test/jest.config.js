module.exports = {
  testEnvironment: "node",
  transform: {
    "^.+\\.ts$": "ts-jest",
  },
  moduleFileExtensions: ["ts", "js"],
  testRegex: "(/specs/.*|(\\.|/)spec)\\.(ts|js)$",
};
