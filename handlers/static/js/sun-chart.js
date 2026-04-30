// Simple canvas-based chart for sun path visualization
function formatHourLabel(hour) {
    if (hour === 0 || hour === 24) {
        return '12AM';
    }
    if (hour === 12) {
        return '12PM';
    }
    if (hour < 12) {
        return `${hour}AM`;
    }
    return `${hour - 12}PM`;
}

function parseSelectedTimeToDecimal(selectedTime) {
    if (!selectedTime) {
        return null;
    }

    const trimmed = selectedTime.trim().toUpperCase();
    const match = trimmed.match(/^([0-9]{1,2}):([0-9]{2})(?::([0-9]{2}))?\s*(AM|PM)?$/);
    if (!match) {
        return null;
    }

    let hour = Number(match[1]);
    const minute = Number(match[2]);
    const second = match[3] ? Number(match[3]) : 0;
    const meridiem = match[4];

    if (meridiem) {
        if (hour === 12) {
            hour = meridiem === 'AM' ? 0 : 12;
        } else if (meridiem === 'PM') {
            hour += 12;
        }
    }

    return hour + minute / 60.0 + second / 3600.0;
}

function drawSunPathChart(ctx, sunPositions, selectedTime = null, peakTime = null) {
    // Clear canvas
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    // Draw sky gradient background
    const gradient = ctx.createLinearGradient(0, 0, 0, ctx.canvas.height);
    gradient.addColorStop(0, '#87CEEB'); // Sky blue at top
    gradient.addColorStop(1, '#FFA07A'); // Light salmon at bottom (for sunrise/sunset)
    ctx.fillStyle = gradient;
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    // Filter out positions with null altitude values
    const validPositions = sunPositions.filter(pos => pos.altitude !== null);

    // Draw night time regions (altitude < -18 degrees, astronomical twilight)
    if (validPositions.length > 1) {
        let inNightRegion = false;
        let nightStartX = 0;

        for (let i = 0; i < validPositions.length; i++) {
            const pos = validPositions[i];
            const x = mapRange(pos.hour, 0, 24, 0, ctx.canvas.width);

            if (pos.altitude < -18) {
                if (!inNightRegion) {
                    // Start of night region
                    nightStartX = x;
                    inNightRegion = true;
                }
            } else if (inNightRegion) {
                // End of night region
                // Draw dark gradient for night
                const nightGradient = ctx.createLinearGradient(nightStartX, 0, x, 0);
                nightGradient.addColorStop(0, 'rgba(0, 0, 40, 0.6)');
                nightGradient.addColorStop(0.5, 'rgba(0, 0, 60, 0.7)');
                nightGradient.addColorStop(1, 'rgba(0, 0, 40, 0.6)');
                ctx.fillStyle = nightGradient;
                ctx.fillRect(nightStartX, 0, x - nightStartX, ctx.canvas.height);

                // Add night label
                ctx.fillStyle = 'rgba(200, 200, 255, 0.8)';
                ctx.font = 'bold 11px Arial';
                const centerX = (nightStartX + x) / 2;
                ctx.fillText('🌙 Night', centerX - 25, 20);

                inNightRegion = false;
            }
        }

        // If still in night region at the end
        if (inNightRegion) {
            const nightGradient = ctx.createLinearGradient(nightStartX, 0, ctx.canvas.width, 0);
            nightGradient.addColorStop(0, 'rgba(0, 0, 40, 0.6)');
            nightGradient.addColorStop(1, 'rgba(0, 0, 60, 0.7)');
            ctx.fillStyle = nightGradient;
            ctx.fillRect(nightStartX, 0, ctx.canvas.width - nightStartX, ctx.canvas.height);

            ctx.fillStyle = 'rgba(200, 200, 255, 0.8)';
            ctx.font = 'bold 11px Arial';
            const centerX = (nightStartX + ctx.canvas.width) / 2;
            ctx.fillText('🌙 Night', centerX - 25, 20);
        }
    }

    // Draw horizon line
    ctx.beginPath();
    ctx.moveTo(0, ctx.canvas.height / 2);
    ctx.lineTo(ctx.canvas.width, ctx.canvas.height / 2);
    ctx.strokeStyle = '#555';
    ctx.lineWidth = 1;
    ctx.stroke();

    // Draw sun path if we have data
    if (validPositions.length > 1) {
        ctx.beginPath();

        // Move to first valid point
        const firstPoint = validPositions[0];
        const x1 = mapRange(firstPoint.hour, 0, 24, 0, ctx.canvas.width);
        const y1 = mapRange(firstPoint.altitude, -90, 90, ctx.canvas.height, 0);
        ctx.moveTo(x1, y1);

        // Draw line through all valid points
        for (let i = 1; i < validPositions.length; i++) {
            const point = validPositions[i];
            const x = mapRange(point.hour, 0, 24, 0, ctx.canvas.width);
            const y = mapRange(point.altitude, -90, 90, ctx.canvas.height, 0);
            ctx.lineTo(x, y);
        }

        ctx.strokeStyle = '#FFA500'; // Orange for sun path
        ctx.lineWidth = 3;
        ctx.stroke();

        // Draw max altitude marker and label
        const maxPosition = validPositions.reduce((best, item) => {
            if (!best) {
                return item;
            }
            return item.altitude > best.altitude ? item : best;
        }, validPositions[0]);

        if (maxPosition && maxPosition.altitude !== null) {
            const maxX = mapRange(maxPosition.hour, 0, 24, 0, ctx.canvas.width);
            const maxY = mapRange(maxPosition.altitude, -90, 90, ctx.canvas.height, 0);

            ctx.save();
            ctx.strokeStyle = 'rgba(255, 255, 255, 0.7)';
            ctx.setLineDash([4, 4]);
            ctx.beginPath();
            ctx.moveTo(maxX, 0);
            ctx.lineTo(maxX, ctx.canvas.height);
            ctx.stroke();
            ctx.restore();

            const text = `Max ${maxPosition.altitude.toFixed(1)}°`;
            const labelX = Math.min(maxX + 8, ctx.canvas.width - 120);
            const labelY = Math.max(20, maxY - 12);
            ctx.fillStyle = 'rgba(255, 255, 255, 0.9)';
            ctx.fillRect(labelX - 4, labelY - 16, 120, 20);
            ctx.fillStyle = '#333';
            ctx.font = '11px Arial';
            ctx.fillText(text, labelX, labelY);
        }

        // Draw peak position marker if available
        if (peakTime) {
            const peakDecimalHour = parseSelectedTimeToDecimal(peakTime);
            if (peakDecimalHour !== null) {
                const peakX = mapRange(peakDecimalHour, 0, 24, 0, ctx.canvas.width);
                ctx.save();
                ctx.strokeStyle = 'rgba(255, 255, 255, 0.9)';
                ctx.setLineDash([4, 4]);
                ctx.beginPath();
                ctx.moveTo(peakX, 0);
                ctx.lineTo(peakX, ctx.canvas.height);
                ctx.stroke();
                ctx.restore();

                ctx.fillStyle = 'rgba(255,255,255,0.9)';
                ctx.fillRect(peakX - 28, ctx.canvas.height - 28, 56, 18);
                ctx.fillStyle = '#333';
                ctx.font = '11px Arial';
                ctx.fillText('Peak', peakX - 20, ctx.canvas.height - 14);
            }
        }

        // Draw sun at current position based on the selected time
        if (validPositions.length > 0) {
            let selectedDecimalHour = null;
            if (selectedTime) {
                selectedDecimalHour = parseSelectedTimeToDecimal(selectedTime);
            }

            let currentPosition;
            if (selectedDecimalHour !== null) {
                currentPosition = validPositions.reduce((best, item) => {
                    return Math.abs(item.hour - selectedDecimalHour) < Math.abs(best.hour - selectedDecimalHour) ? item : best;
                }, validPositions[0]);
            } else {
                currentPosition = validPositions[validPositions.length - 1];
            }

            const x = mapRange(currentPosition.hour, 0, 24, 0, ctx.canvas.width);
            const y = mapRange(currentPosition.altitude, -90, 90, ctx.canvas.height, 0);

            // Draw sun
            ctx.beginPath();
            ctx.arc(x, y, 12, 0, Math.PI * 2);
            ctx.fillStyle = '#FFD700'; // Gold color for sun
            ctx.fill();
            ctx.strokeStyle = '#FFA500';
            ctx.lineWidth = 2;
            ctx.stroke();

            // Draw glow effect
            ctx.beginPath();
            ctx.arc(x, y, 18, 0, Math.PI * 2);
            ctx.strokeStyle = 'rgba(255, 215, 0, 0.6)';
            ctx.lineWidth = 3;
            ctx.stroke();

            // Draw sun rays
            ctx.strokeStyle = 'rgba(255, 165, 0, 0.5)';
            ctx.lineWidth = 2;
            for (let i = 0; i < 8; i++) {
                const angle = (i * Math.PI / 4);
                const rayX1 = x + Math.cos(angle) * 15;
                const rayY1 = y + Math.sin(angle) * 15;
                const rayX2 = x + Math.cos(angle) * 22;
                const rayY2 = y + Math.sin(angle) * 22;

                ctx.beginPath();
                ctx.moveTo(rayX1, rayY1);
                ctx.lineTo(rayX2, rayY2);
                ctx.stroke();
            }
        }
    }
    
    // Draw time markers on x-axis
    ctx.fillStyle = '#333';
    ctx.font = '10px Arial';
    for (let hour = 0; hour <= 24; hour += 3) {
        const x = mapRange(hour, 0, 24, 0, ctx.canvas.width);
        ctx.fillText(formatHourLabel(hour), x - 15, ctx.canvas.height - 5);
        
        // Draw vertical grid line
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, ctx.canvas.height);
        ctx.strokeStyle = 'rgba(0, 0, 0, 0.1)';
        ctx.lineWidth = 1;
        ctx.stroke();
    }
    
    // Draw altitude markers on y-axis
    ctx.fillStyle = '#333';
    ctx.font = '10px Arial';
    for (let alt = -90; alt <= 90; alt += 30) {
        const y = mapRange(alt, -90, 90, ctx.canvas.height, 0);
        ctx.fillText(`${alt}°`, 5, y + 3);
        
        // Draw horizontal grid line
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(ctx.canvas.width, y);
        ctx.strokeStyle = 'rgba(0, 0, 0, 0.1)';
        ctx.lineWidth = 1;
        ctx.stroke();
    }
}


// Function to calculate sun positions for the entire day
async function calculateDailySunPath(lat, lon, date, cityName = null, selectedTime = null) {
    const positions = [];

    let baseUrl;
    if (cityName) {
        baseUrl = `/sun-pos/api/sun-position?city=${encodeURIComponent(cityName)}&date=${date}`;
    } else {
        baseUrl = `/sun-pos/api/sun-position?lat=${lat}&lon=${lon}&date=${date}`;
    }

    for (let hour = 0; hour < 24; hour++) {
        for (let minute = 0; minute < 60; minute += 15) { // Every 15 minutes
            const timeStr = `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}`;
            try {
                const response = await fetch(`${baseUrl}&time=${timeStr}`);
                const data = await response.json();

                positions.push({
                    hour: hour + minute/60,
                    altitude: data.sun_altitude,
                    azimuth: data.sun_azimuth,
                    time: timeStr
                });
            } catch (error) {
                console.error(`Error calculating sun position for ${timeStr}:`, error);
            }
        }
    }

    return positions;
}

// Helper function to map a value from one range to another
function mapRange(value, inMin, inMax, outMin, outMax) {
    return (value - inMin) * (outMax - outMin) / (inMax - inMin) + outMin;
}