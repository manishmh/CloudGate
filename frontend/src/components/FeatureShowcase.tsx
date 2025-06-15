"use client";

import { useState } from "react";
import { HiChevronLeft, HiChevronRight } from "react-icons/hi";

interface Feature {
  id: string;
  title: string;
  description: string;
  icon: string;
  stats: string;
  color: string;
  features: readonly string[];
}

interface FeatureShowcaseProps {
  features: readonly Feature[];
  title?: string;
  subtitle?: string;
}

export default function FeatureShowcase({
  features,
  title = "Platform Features",
  subtitle = "Enterprise-grade security and performance",
}: FeatureShowcaseProps) {
  const [currentIndex, setCurrentIndex] = useState(0);

  const nextFeature = () => {
    setCurrentIndex((prev) => (prev + 1) % features.length);
  };

  const prevFeature = () => {
    setCurrentIndex((prev) => (prev - 1 + features.length) % features.length);
  };

  const goToFeature = (index: number) => {
    setCurrentIndex(index);
  };

  if (!features.length) return null;

  return (
    <div className="relative">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <h3 className="text-2xl font-bold text-gray-900">{title}</h3>
          <p className="text-gray-600 mt-1">{subtitle}</p>
        </div>

        {/* Navigation Controls */}
        <div className="hidden md:flex items-center space-x-2">
          <button
            onClick={prevFeature}
            className="p-2 rounded-xl bg-gray-100 hover:bg-gray-200 transition-colors cursor-pointer"
            aria-label="Previous feature"
          >
            <HiChevronLeft className="h-5 w-5 text-gray-600" />
          </button>
          <button
            onClick={nextFeature}
            className="p-2 rounded-xl bg-gray-100 hover:bg-gray-200 transition-colors cursor-pointer"
            aria-label="Next feature"
          >
            <HiChevronRight className="h-5 w-5 text-gray-600" />
          </button>
        </div>
      </div>

      {/* Desktop Grid View */}
      <div className="hidden lg:grid grid-cols-4 gap-8">
        {features.map((feature) => (
          <div
            key={feature.id}
            className="group relative bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-xl hover:border-gray-200 transition-all duration-300 cursor-pointer overflow-hidden"
          >
            {/* Gradient Background */}
            <div
              className={`absolute inset-0 bg-gradient-to-br ${feature.color} opacity-0 group-hover:opacity-5 transition-opacity duration-300`}
            />

            <div className="relative z-10">
              <div className="flex items-center justify-between mb-6">
                <div className="text-4xl transform group-hover:scale-110 transition-transform duration-200">
                  {feature.icon}
                </div>
                <span className="text-xs font-semibold text-gray-500 bg-gray-100 px-3 py-1.5 rounded-full group-hover:bg-white group-hover:text-gray-700 transition-colors">
                  {feature.stats}
                </span>
              </div>

              <h4 className="text-xl font-semibold text-gray-900 mb-3 group-hover:text-gray-800">
                {feature.title}
              </h4>
              <p className="text-gray-600 text-sm mb-6 leading-relaxed group-hover:text-gray-700">
                {feature.description}
              </p>

              <div className="space-y-2">
                {feature.features.map((item, featureIndex) => (
                  <div
                    key={featureIndex}
                    className="flex items-center text-sm text-gray-500 group-hover:text-gray-600 transition-colors"
                  >
                    <div className="w-1.5 h-1.5 bg-gray-400 rounded-full mr-3 group-hover:bg-gray-500" />
                    {item}
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Mobile/Tablet Carousel View */}
      <div className="lg:hidden">
        {/* Feature Cards Container */}
        <div className="relative overflow-hidden rounded-2xl">
          <div
            className="flex transition-transform duration-500 ease-in-out"
            style={{ transform: `translateX(-${currentIndex * 100}%)` }}
          >
            {features.map((feature) => (
              <div key={feature.id} className="w-full flex-shrink-0 p-6">
                <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 h-full">
                  {/* Gradient Background */}
                  <div
                    className={`absolute inset-0 bg-gradient-to-br ${feature.color} opacity-5 rounded-2xl`}
                  />

                  <div className="relative z-10 text-center">
                    <div className="text-6xl mb-6">{feature.icon}</div>
                    <span className="inline-block text-sm font-semibold text-gray-500 bg-gray-100 px-3 py-1 rounded-full mb-4">
                      {feature.stats}
                    </span>

                    <h4 className="text-2xl font-bold text-gray-900 mb-4">
                      {feature.title}
                    </h4>
                    <p className="text-gray-600 mb-6 text-lg">
                      {feature.description}
                    </p>

                    <div className="space-y-2">
                      {feature.features.map((item, featureIndex) => (
                        <div
                          key={featureIndex}
                          className="flex items-center justify-center text-sm text-gray-500"
                        >
                          <div className="w-1.5 h-1.5 bg-gray-400 rounded-full mr-3" />
                          {item}
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Dots Indicator */}
        <div className="flex justify-center mt-6 space-x-2">
          {features.map((_, index) => (
            <button
              key={index}
              onClick={() => goToFeature(index)}
              className={`w-2 h-2 rounded-full transition-all duration-200 cursor-pointer ${
                index === currentIndex
                  ? "bg-blue-600 w-8"
                  : "bg-gray-300 hover:bg-gray-400"
              }`}
              aria-label={`Go to feature ${index + 1}`}
            />
          ))}
        </div>

        {/* Mobile Navigation */}
        <div className="flex justify-center mt-4 space-x-4 md:hidden">
          <button
            onClick={prevFeature}
            className="flex items-center px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-xl transition-colors cursor-pointer"
          >
            <HiChevronLeft className="h-4 w-4 mr-1" />
            Previous
          </button>
          <button
            onClick={nextFeature}
            className="flex items-center px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-xl transition-colors cursor-pointer"
          >
            Next
            <HiChevronRight className="h-4 w-4 ml-1" />
          </button>
        </div>
      </div>

      {/* Feature Counter */}
      <div className="text-center mt-8 lg:hidden">
        <span className="text-sm text-gray-500">
          {currentIndex + 1} of {features.length}
        </span>
      </div>
    </div>
  );
}
