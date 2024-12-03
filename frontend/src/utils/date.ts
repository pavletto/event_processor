export const formatDate = (timestamp?: number | Date, pattern = '{Y}-{M}-{D} {H}:{m}:{s}', utc = true): string => {
    if (!timestamp) {
        return "";
    }
    if (timestamp instanceof Date) {
        timestamp = timestamp.getTime();
    }
    const d = new Date(timestamp);
    const Y = utc ? d.getUTCFullYear() : d.getFullYear(),
        M = String((utc ? d.getUTCMonth() : d.getMonth()) + 1).padStart(2, '0'),
        D = String(utc ? d.getUTCDate() : d.getDate()).padStart(2, '0'),
        hours = utc ? d.getUTCHours() : d.getHours(),
        is12HourFormat = pattern.includes('{A}'),
        H = is12HourFormat
            ? String(hours % 12 || 12).padStart(2, '0') // 12-hour format if '{A}' is in the pattern
            : String(hours).padStart(2, '0'),           // 24-hour format otherwise
        m = String(utc ? d.getUTCMinutes() : d.getMinutes()).padStart(2, '0'),
        s = String(utc ? d.getUTCSeconds() : d.getSeconds()).padStart(2, '0'),
        ms = String(utc ? d.getUTCMilliseconds() : d.getMilliseconds()).padStart(3, '0'),
        A = is12HourFormat && hours < 12 ? 'AM' : is12HourFormat ? 'PM' : ''; // Set 'AM/PM' only if in 12-hour format

    return stringTemplate(pattern, { Y, M, D, H, m, s, ms, A });
};

export function stringTemplate(template: string, params?: any): string {
    return template.replace(/{[^{}]+}/g, function (key) {
        return params[key.replace(/[{}]+/g, "")] || "";
    });
}