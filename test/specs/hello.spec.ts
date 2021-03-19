import api from "../helpers/api";

describe("hello endpoint test", () => {
  test("success", async () => {
    const res = await api.get("/");
    expect(res.status).toEqual(200);
    expect(res.data.length).toBeGreaterThan(0);
  });
});
