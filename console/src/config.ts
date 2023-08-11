// const apiEndpoint = process.env.API_ENDPOINT || "/api";
console.log(import.meta.env);
const apiEndpoint = import.meta.env.VITE_API_ENDPOINT || "/api";

export { apiEndpoint };
