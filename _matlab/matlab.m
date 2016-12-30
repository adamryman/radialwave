%%% Original maplab code, new comments start with %%%
%%% Done in University of Washington EE 341 Continuous Signal Processing
%%% By Adam Ryman, Unlicensed
clear all
close all

%extract max freqs

%%% 10 breaks
%%% The more breaks the more chords, but it depends on the sampling rate, the length of the song to determine the chords
%%% musicLength / (SamplingRate / breaks) = number of chords
breaks = 10;

%%% Wavelenght of light to rgb color, starting at violet/blue going to red
mycolors = csvread('wavToRGB.csv');

%%% Ok music and sample rate (N)
%%% harder.wav would be the file to process
[music, N] = wavread('harder.wav');

%%% Length of the music
mySize = size(music,1);

%%%  same as [0, 1, 2, 3, ..., N-1]
%%%  (0:(N-1))
%%% It seems that evey element is multipled by N and then divied by N which should be the original array
%%% N*(Array)/N = (N/N)*Array = 1*Array = Array
x = N*(0:(N-1))/N;
%%% x seems unused after this, only got 1 x2 after
%%% Subtract every element by N/2
%%% Should be something like [-N/2, 1-N/2, 2-N/2, ... (N-1)-N/2] Possibly a zero in the middle
%%% These are frequencies from -N/2 to 0 to N/2, basically what we will get after an FFT shift
%%% If we do not FFT shift, I think we use x, but there may be an issue with that
x2 = x - N/2;

%%% Just allocating an array of the right size
%%% golang
%%% frequencies = make([]float64, mySize/(N/breaks))
%%% Need to round up though, like the ceil func
frequencies = zeros(1,(ceil(mySize/(N/breaks))));

%%% A position in frequencies
q = 1;

%%% Break the sample rate and advance by that amount until we hit the end of the song
%%% Probably should have started at 0 (or in matlab, 1, because fuck you 0 indexing)
%%% Ahhh this is why it looks like we miss the last one, because we do
for n = N/breaks:N/breaks:mySize

	%%% Piece by piece
	%%%
	%%% Take the spot we are now, go back by that amount but advance by 1
	%%% n - N/breaks+ 1 >>> startn
	%%%
	%%% Take the music from startn to n
	%%% music( startn : n) >>> musicSlice
	%%%
	%%% fft that with the original sampling rate (N)
	%%% ya = fft(musicSlice, N)
	%%%
    ya = fft(music(n - N/breaks+ 1 : n), N);
	%%%
	%%% This shifts the zero-frequency component to center of the array
	%%% We are using x2, which has the zero-frequency component in the center
	%%%
    ya2 = fftshift(ya);
	%%%
	%%% First, positive and negative values are the same
	%%% Find the largest value in the y axis (value = amplitude; index = frequency)
	%%% Give us the value and the posiiton
	%%%
	%%% _, i = max(abs(ya2))
	%%%
    [a,i] = max(abs(ya2));
	%%%
	%%% q is a counter, we are using it to append new values
	%%% in golang: frequencies = append(frequencies, newValue)
	%%% abs(x2(i)) >>> newValue
	%%% What is ab(x2(i))
	%%%
	%%% The values of x2 are the frequency values for each index given to us from ya2
	%%% I believe that the values of x relate to ya in the same way
	%%%
	%%% We store this frequency value (the dominate freuence for this time slice) in frequencies
	%%%
    frequencies(q) =  abs(x2(i));
    q = q + 1;
	%%%
	%%%
	%%% I think this should have had the form
	%%% n = 1:N/breaks:mySize
	%%% Or
	%%% var n int
	%%% for n = 0; n < mySize; n = n + N/breaks {
	%%%  music(n : n + N/breaks)
	%%%  ...
	%%%	}
	%%% We would want to do the last section too
	%%% n - N/breaks to mySize

	%%%

end

%%% Remove duplicates, we are creating a smooth mapping of frequencies to colors
%%% Also sorts them, so u will create our mapping
u = unique(frequencies);

%%% Create length(frequencies) partitions from 0 to 2*pi
angle = linspace(0, 2*pi, length(frequencies));

%%% angle to x pos chord end
c = cos(angle);
%%% angle to y pos chord end
d = sin(angle);

hold on;
box off;
set(gca,'YTick',[]);
set(gca,'XTick',[]);
set(gca,'xcolor',[1 1 1]);
set(gca,'ycolor',[1 1 1]);
plot(c,d);

for i = 1:length(frequencies)
	%%% Find the index of u that has the same value as the current frequency.
	%%% Take length(myColors)=266 ; length(uniqueFreq) and divide them 266 /
	%%% length(uniqueFreq) so that each index of uniqueFreq will represent equal
	%%% distance though the color space, but if uniqueFreq < 266 then we jump
	%%% around in the space, skipping gaps. This mapping could be better.
    freqColor = mycolors(ceil(find(u==frequencies(i)).* length(mycolors) ./ length(u)) , :);
	%%% All chords start at 0,0 and go to x(i), y(i)
    plot([0,c(i)],[0,d(i)], 'Color', freqColor);
    axis('equal');
end

F = getframe();

original = F.cdata;

imwrite(original, 'beforeEdge.png', 'png');

% onto edge detection
image = rgb2gray(original);

h1 = [-1 0 1; -2 0 2; -1 0 1];

h2 = [1 2 1; 0 0 0; -1 -2 -1];

m1 = conv2(h1,image);

m2 = conv2(h2,image);

mn = (m1.^2 + m2.^2).^0.5;

m1 = uint8(m1);

m2 = uint8(m2);

mn = uint8(mn);

%lets multiple that edge dection throught the created image
someExtraOnes = [original ones(size(original,1),2,3)];
extraOnes = [someExtraOnes; ones(2,size(someExtraOnes,2),3)];

a = extraOnes(:,:,1) .* mn;
b = extraOnes(:,:,2) .* mn;
c = extraOnes(:,:,3) .* mn;

d = cat(3,a,b,c);

figure

imshow(d)
imwrite(d, 'harder.png', 'png');
