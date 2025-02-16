export function randomId() {
    const prefix = String(Math.random());
    if (prefix.length == 20)
        return prefix; 
    return prefix + String(Math.random()).substring(0, 20 - prefix.length)
}