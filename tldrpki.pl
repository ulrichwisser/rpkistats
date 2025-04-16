#!/usr/bin/perl

use strict;
use warnings;

my $TLD = $ARGV[0];
$TLD .= '.' if substr($TLD, -1) ne '.';

my $nsraw = `dig \@9.9.9.9 $TLD NS +noall +answer`;

my @nsraw = split /\n/, $nsraw;

my %ip4s = ();
my %ip6s = ();

my %name2ip4 = ();
my %name2ip6 = ();

my @ns = ();

for my $nsrr (@nsraw) {
    $nsrr =~ m/^\S+\s+\d+\s+IN\s+NS\s+(\S+)\s*$/;
    my $ns = $1;
    push @ns, $ns;
    $name2ip4{$ns} = [];
    $name2ip6{$ns} = [];
    my $ip4raw = `dig \@9.9.9.9 $ns A +noall +answer`;
    my @ip4raw = split /\n/, $ip4raw;
    for my $iprr (@ip4raw) {
        $iprr =~ m/^\S+\s+\d+\s+IN\s+A\s+(\S+)$/;
        my $ip = $1;
        $ip4s{$ip} = 0;
        push @{$name2ip4{$ns}}, $ip;
    } 
    my $ip6raw = `dig \@9.9.9.9 $ns AAAA +noall +answer`;
    my @ip6raw = split /\n/, $ip6raw;
    for my $iprr (@ip6raw) {
        $iprr =~ m/^\S+\s+\d+\s+IN\s+AAAA\s+(\S+)$/;
        my $ip = $1;
        $ip6s{$ip} = 0;
        push @{$name2ip6{$ns}}, $ip;
    } 
}

my %tas4 = ();
my %as4 = ();

for my $ip4 (keys %ip4s) {
    $ip4 =~ m/(\d+\.\d+\.\d+\.)\d+/;
    my $prefix = $1 . '0/24';
    my $roaraw = `curl -s http://void.home.wisser.se:8323/json?select-prefix=$prefix`;
    my @roalines = split /\n/, $roaraw;
    foreach my $line (@roalines) {
        next if !($line =~ m/{ "asn": "(AS\d+)", "prefix": "\S+", "maxLength": \d+, "ta": "(\S+)" }/);
        $as4{$1} = 1;
        $tas4{$2} = 1;
        $ip4s{$ip4} = 1;
        print "prefix $prefix from $2 $1\n";
    }
}

my %tas6 = ();
my %as6 = ();

for my $ip6 (keys %ip6s) {
    my @v6parts = split /:/, $ip6;
    my $prefixlen = 0;
    my $prefix = "";
    foreach my $part (@v6parts) {
        last if $part eq "";
        $prefix .= ':' if length($prefix) > 0;
        $prefix .= $part;
        $prefixlen += 16;
        last if $prefixlen == 64;
    }
    $prefix .= "::/$prefixlen";
    my $roaraw = `curl -s http://void.home.wisser.se:8323/json?select-prefix=$prefix`;
    my @roalines = split /\n/, $roaraw;
    foreach my $line (@roalines) {
        next if !($line =~ m/{ "asn": "(AS\d+)", "prefix": "\S+", "maxLength": \d+, "ta": "(\S+)" }/);
        $as6{$1} = 1;
        $tas6{$2} = 1;
        $ip6s{$ip6} = 1;
        print "prefix $prefix from $2 $1\n";
    }
}

my $names_full = 0;
my $names_partial = 0;

foreach my $ns (@ns) {
    my $roas = 0;
    foreach my $ip4 (@{$name2ip4{$ns}}) {
        $roas++ if $ip4s{$ip4};
    }
    foreach my $ip6 (@{$name2ip6{$ns}}) {
        $roas++ if $ip6s{$ip6};
    }
    if ($roas == scalar(@{$name2ip4{$ns}})+scalar(@{$name2ip6{$ns}})) {
        print "$ns is full\n";
        $names_full++;
    } elsif ($roas > 0) {
        print "$ns is partial\n";
        $names_partial++;
    }
}

my $ip4s_roas = 0;
foreach my $ip4 (keys %ip4s) {
    $ip4s_roas++ if $ip4s{$ip4};
}

my $ip6s_roas = 0;
foreach my $ip6 (keys %ip6s) {
    $ip6s_roas++ if $ip6s{$ip6};
}

my $names = scalar(@nsraw);
my $ip4s  = scalar(keys %ip4s);
my $ip6s  = scalar(keys %ip6s);
my $tas4  = scalar(keys %tas4);
my $tas6  = scalar(keys %tas6);
my $as4   = scalar(keys %as4);
my $as6   = scalar(keys %as6);

print "NAMES:         $names\n";
print "NAMES FULL     $names_full\n";
print "NAMES PARTIAL  $names_partial\n";
print "IP4:           $ip4s\n";
print "IP4 ROAS:      $ip4s_roas\n";
print "IP6:           $ip6s\n";
print "IP6 ROAS:      $ip6s_roas\n";
print "TAS4:          $tas4\n";
print "TAS6:          $tas6\n";
print "AS4:           $as4\n";
print "AS6:           $as6\n";
