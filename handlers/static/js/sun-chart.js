// Simple canvas-based chart for sun path visualization
function drawSunPathChart(ctx, sunPositions, selectedTime = null) {
    // Clear canvas
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
    
    // Draw sky gradient background
    const gradient = ctx.createLinearGradient(0, 0, 0, ctx.canvas.height);
    gradient.addColorStop(0, '#87CEEB'); // Sky blue at top
    gradient.addColorStop(1, '#FFA07A'); // Light salmon at bottom (for sunrise/sunset)
    ctx.fillStyle = gradient;
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height);
    
    // Draw horizon line
    ctx.beginPath();
    ctx.moveTo(0, ctx.canvas.height / 2);
    ctx.lineTo(ctx.canvas.width, ctx.canvas.height / 2);
    ctx.strokeStyle = '#555';
    ctx.lineWidth = 1;
    ctx.stroke();
    
    // Draw sun path if we have data
    // Filter out positions with null altitude values
    const validPositions = sunPositions.filter(pos => pos.altitude !== null);

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
        
        // Draw sun at current position based on the selected time
        if (validPositions.length > 0) {
            // Use the selectedTime parameter passed from the calling function
            if (selectedTime) {
                // Parse the selected time (HH:MM format)
                const [hours, minutes] = selectedTime.split(':').map(Number);
                const selectedDecimalHour = hours + minutes / 60.0;

                // Find the closest valid position to the selected time
                let closestPosition = validPositions[0];
                let minDiff = Math.abs(closestPosition.hour - selectedDecimalHour);

                for (let i = 1; i < validPositions.length; i++) {
                    const diff = Math.abs(validPositions[i].hour - selectedDecimalHour);
                    if (diff < minDiff) {
                        minDiff = diff;
                        closestPosition = validPositions[i];
                    }
                }

                const x = mapRange(closestPosition.hour, 0, 24, 0, ctx.canvas.width);
                const y = mapRange(closestPosition.altitude, -90, 90, ctx.canvas.height, 0);

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
            } else {
                // Fallback to drawing at the last valid position if no specific time is provided
                const current = validPositions[validPositions.length - 1]; // Last valid position
                const x = mapRange(current.hour, 0, 24, 0, ctx.canvas.width);
                const y = mapRange(current.altitude, -90, 90, ctx.canvas.height, 0);

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
    }
    
    // Draw time markers on x-axis
    ctx.fillStyle = '#333';
    ctx.font = '10px Arial';
    for (let hour = 0; hour <= 24; hour += 3) {
        const x = mapRange(hour, 0, 24, 0, ctx.canvas.width);
        ctx.fillText(`${hour}:00`, x - 10, ctx.canvas.height - 5);
        
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