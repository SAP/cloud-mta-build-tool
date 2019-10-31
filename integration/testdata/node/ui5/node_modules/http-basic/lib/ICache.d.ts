import { CachedResponse } from './CachedResponse';
interface ICache {
    getResponse(url: string, cb: (err: Error | null, response: CachedResponse | null) => void): void;
    setResponse(url: string, response: CachedResponse | null): void;
    invalidateResponse(url: string, cb: (err: Error | null) => void): void;
}
export { ICache };
