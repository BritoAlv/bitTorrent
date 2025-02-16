export function randomId() {
    const prefix = String(Math.random()).substring(2);
    if (prefix.length == 20)
        return prefix; 
    return prefix + String(Math.random()).substring(2, 22 - prefix.length)
}