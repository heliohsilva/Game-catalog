'use client';

import React, { useState } from 'react';
import { 
  X, 
  Plus, 
  Trash2, 
  Layers, 
  AlertCircle, 
  Loader2, 
  CheckCircle2, 
  Tag
} from 'lucide-react';
import { PlatformInfo } from '../types/game';

interface ManagePlatformsModalProps {
  isOpen: boolean;
  onClose: () => void;
  platforms: PlatformInfo[];
  onPlatformsChanged: () => void;
}

export const ManagePlatformsModal: React.FC<ManagePlatformsModalProps> = ({
  isOpen,
  onClose,
  platforms,
  onPlatformsChanged,
}) => {
  // New platform form state
  const [newPlatformName, setNewPlatformName] = useState('');
  const [newPlatformSubs, setNewPlatformSubs] = useState('');
  const [isAddingPlatform, setIsAddingPlatform] = useState(false);

  // New subcategory per-platform state
  const [addingSubFor, setAddingSubFor] = useState<string | null>(null);
  const [newSubName, setNewSubName] = useState('');
  const [isAddingSub, setIsAddingSub] = useState(false);

  // Status message
  const [feedback, setFeedback] = useState<{ type: 'success' | 'error'; message: string } | null>(null);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  if (!isOpen) return null;

  const showFeedback = (type: 'success' | 'error', message: string) => {
    setFeedback({ type, message });
    setTimeout(() => {
      setFeedback(null);
    }, 4000);
  };

  // Add a new platform
  const handleAddPlatform = async (e: React.FormEvent) => {
    e.preventDefault();
    const name = newPlatformName.trim();
    if (!name) return;

    setIsAddingPlatform(true);
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

    const subcategories = newPlatformSubs
      .split(',')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    try {
      const res = await fetch(`${apiUrl}/platforms`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, subcategories }),
      });

      const data = await res.json();
      if (!res.ok) {
        showFeedback('error', data.message || 'Failed to create platform');
        return;
      }

      setNewPlatformName('');
      setNewPlatformSubs('');
      showFeedback('success', `Platform '${name}' created successfully`);
      onPlatformsChanged();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Network error';
      showFeedback('error', msg);
    } finally {
      setIsAddingPlatform(false);
    }
  };

  // Delete a platform
  const handleDeletePlatform = async (platformName: string, gameCount: number = 0) => {
    const confirmMsg = gameCount > 0
      ? `Platform '${platformName}' contains ${gameCount} game(s). Deleting it will also remove these games. Are you sure?`
      : `Are you sure you want to remove platform '${platformName}'?`;

    if (!window.confirm(confirmMsg)) {
      return;
    }

    setActionLoading(`delete-plat-${platformName}`);
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

    try {
      const res = await fetch(`${apiUrl}/platforms/${encodeURIComponent(platformName)}?force=true`, {
        method: 'DELETE',
      });
      const data = await res.json();
      if (!res.ok) {
        showFeedback('error', data.message || 'Failed to delete platform');
        return;
      }

      showFeedback('success', `Platform '${platformName}' removed`);
      onPlatformsChanged();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Network error';
      showFeedback('error', msg);
    } finally {
      setActionLoading(null);
    }
  };

  // Add a subcategory to an existing platform
  const handleAddSubcategory = async (platformName: string) => {
    const sub = newSubName.trim();
    if (!sub) return;

    setIsAddingSub(true);
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

    try {
      const res = await fetch(`${apiUrl}/platforms/${encodeURIComponent(platformName)}/subcategories`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: sub }),
      });
      const data = await res.json();
      if (!res.ok) {
        showFeedback('error', data.message || 'Failed to add subcategory');
        return;
      }

      setNewSubName('');
      setAddingSubFor(null);
      showFeedback('success', `Subplatform '${sub}' added to ${platformName}`);
      onPlatformsChanged();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Network error';
      showFeedback('error', msg);
    } finally {
      setIsAddingSub(false);
    }
  };

  // Delete a subcategory from a platform
  const handleDeleteSubcategory = async (platformName: string, subcategory: string) => {
    if (!window.confirm(`Remove subplatform '${subcategory}' from ${platformName}?`)) {
      return;
    }

    setActionLoading(`delete-sub-${platformName}-${subcategory}`);
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

    try {
      const res = await fetch(
        `${apiUrl}/platforms/${encodeURIComponent(platformName)}/subcategories/${encodeURIComponent(subcategory)}?force=true`,
        { method: 'DELETE' }
      );
      const data = await res.json();
      if (!res.ok) {
        showFeedback('error', data.message || 'Failed to delete subcategory');
        return;
      }

      showFeedback('success', `Subplatform '${subcategory}' removed from ${platformName}`);
      onPlatformsChanged();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Network error';
      showFeedback('error', msg);
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md"
      role="dialog"
      aria-modal="true"
      aria-label="Manage Platforms and Subplatforms"
    >
      <div 
        className="hypr-glass rounded-2xl w-full max-w-2xl max-h-[85vh] flex flex-col border border-[var(--border-color)] hypr-window-in shadow-2xl overflow-hidden"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--border-color)]">
          <div className="flex items-center gap-2.5">
            <Layers className="w-5 h-5 text-[var(--accent)]" />
            <div>
              <h2 className="font-bold text-sm text-[var(--fg-bright)] tracking-wide">
                Manage Platforms & Subplatforms
              </h2>
              <p className="text-[11px] text-[var(--fg-light)]">
                Add on-demand platforms (e.g. Nintendo DS) and launcher subcategories (e.g. Ubisoft Connect)
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1 rounded-lg text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
            title="Close"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Feedback message banner */}
        {feedback && (
          <div className={`px-5 py-2.5 text-xs flex items-center gap-2 ${
            feedback.type === 'success' 
              ? 'bg-emerald-500/15 text-emerald-300 border-b border-emerald-500/20' 
              : 'bg-rose-500/15 text-rose-300 border-b border-rose-500/20'
          }`}>
            {feedback.type === 'success' ? (
              <CheckCircle2 className="w-3.5 h-3.5 shrink-0" />
            ) : (
              <AlertCircle className="w-3.5 h-3.5 shrink-0" />
            )}
            <span>{feedback.message}</span>
          </div>
        )}

        {/* Body Container */}
        <div className="flex-1 overflow-y-auto p-5 space-y-5">

          {/* Form to add a new Platform */}
          <form 
            onSubmit={handleAddPlatform}
            className="p-3.5 rounded-xl bg-[var(--bg-lighter)]/40 border border-[var(--border-color)] space-y-3"
          >
            <div className="flex items-center gap-2">
              <Plus className="w-3.5 h-3.5 text-[var(--accent)]" />
              <span className="font-bold text-xs text-[var(--fg-bright)]">Add New Platform</span>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <div>
                <label className="block text-[10px] font-mono text-[var(--fg-light)] mb-1">
                  PLATFORM NAME *
                </label>
                <input
                  type="text"
                  placeholder="e.g. Nintendo DS, Dreamcast, Vita"
                  value={newPlatformName}
                  onChange={(e) => setNewPlatformName(e.target.value)}
                  className="w-full px-3 py-1.5 rounded-lg bg-[var(--bg-primary)] border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
                  required
                />
              </div>

              <div>
                <label className="block text-[10px] font-mono text-[var(--fg-light)] mb-1">
                  SUBPLATFORMS (OPTIONAL, COMMA-SEPARATED)
                </label>
                <input
                  type="text"
                  placeholder="e.g. Cartridge, Homebrew, Digital"
                  value={newPlatformSubs}
                  onChange={(e) => setNewPlatformSubs(e.target.value)}
                  className="w-full px-3 py-1.5 rounded-lg bg-[var(--bg-primary)] border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
                />
              </div>
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={isAddingPlatform || !newPlatformName.trim()}
                className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-[var(--accent)] text-[var(--bg-primary)] text-xs font-bold transition-opacity hover:opacity-90 disabled:opacity-50 cursor-pointer"
              >
                {isAddingPlatform ? (
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                ) : (
                  <Plus className="w-3.5 h-3.5" />
                )}
                <span>Add Platform</span>
              </button>
            </div>
          </form>

          {/* Current Platforms List */}
          <div className="space-y-3">
            <h3 className="text-xs font-bold text-[var(--fg-light)] uppercase tracking-wider font-mono">
              Configured Platforms ({platforms.length})
            </h3>

            {platforms.length === 0 ? (
              <div className="text-center py-6 text-xs text-[var(--fg-light)]">
                No platforms configured yet. Add your first platform above.
              </div>
            ) : (
              <div className="space-y-3">
                {platforms.map((p) => {
                  const isDeleting = actionLoading === `delete-plat-${p.name}`;
                  const isAddingSubHere = addingSubFor === p.name;

                  return (
                    <div 
                      key={p.name}
                      className="p-3.5 rounded-xl bg-[var(--bg-card)] border border-[var(--border-color)] space-y-2.5 transition-colors"
                    >
                      {/* Platform header */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="font-bold text-xs text-[var(--fg-bright)]">{p.name}</span>
                          <span className="px-2 py-0.5 rounded-full text-[10px] font-mono bg-[var(--selection)] text-[var(--accent)]">
                            {p.gameCount ?? 0} games
                          </span>
                        </div>

                        <button
                          type="button"
                          onClick={() => handleDeletePlatform(p.name, p.gameCount ?? 0)}
                          disabled={isDeleting}
                          className="flex items-center gap-1 p-1 px-2 rounded-lg text-rose-400 hover:bg-rose-500/15 text-[11px] transition-colors cursor-pointer disabled:opacity-50"
                          title={`Delete platform ${p.name}`}
                        >
                          {isDeleting ? (
                            <Loader2 className="w-3 h-3 animate-spin" />
                          ) : (
                            <Trash2 className="w-3 h-3" />
                          )}
                          <span>Remove</span>
                        </button>
                      </div>

                      {/* Subcategories chips */}
                      <div className="flex flex-wrap items-center gap-1.5 pt-1">
                        <span className="text-[10px] font-mono text-[var(--fg-light)] flex items-center gap-1 mr-1">
                          <Tag className="w-3 h-3" />
                          Subplatforms:
                        </span>

                        {p.subcategories && p.subcategories.length > 0 ? (
                          p.subcategories.map((sub) => {
                            const isDeletingSub = actionLoading === `delete-sub-${p.name}-${sub}`;
                            return (
                              <span
                                key={sub}
                                className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-[var(--bg-lighter)] border border-[var(--border-color)] text-[11px] text-[var(--fg-primary)]"
                              >
                                <span>{sub}</span>
                                <button
                                  type="button"
                                  onClick={() => handleDeleteSubcategory(p.name, sub)}
                                  disabled={isDeletingSub}
                                  className="text-[var(--fg-light)] hover:text-rose-400 cursor-pointer disabled:opacity-50"
                                  title={`Remove ${sub} from ${p.name}`}
                                >
                                  {isDeletingSub ? (
                                    <Loader2 className="w-2.5 h-2.5 animate-spin" />
                                  ) : (
                                    <X className="w-2.5 h-2.5" />
                                  )}
                                </button>
                              </span>
                            );
                          })
                        ) : (
                          <span className="text-[11px] italic text-[var(--fg-light)] opacity-60">
                            None (Direct/Generic)
                          </span>
                        )}

                        {/* Button to toggle add subcategory */}
                        {!isAddingSubHere && (
                          <button
                            type="button"
                            onClick={() => {
                              setAddingSubFor(p.name);
                              setNewSubName('');
                            }}
                            className="inline-flex items-center gap-0.5 px-2 py-0.5 rounded-md border border-dashed border-[var(--accent)]/40 hover:border-[var(--accent)] text-[10px] text-[var(--accent)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
                          >
                            <Plus className="w-2.5 h-2.5" />
                            <span>Add Subplatform</span>
                          </button>
                        )}
                      </div>

                      {/* Inline form to add subcategory to this platform */}
                      {isAddingSubHere && (
                        <div className="flex items-center gap-2 pt-2 border-t border-[var(--border-color)]/50">
                          <input
                            type="text"
                            placeholder="e.g. Ubisoft Connect, EA App, Cartridge"
                            value={newSubName}
                            onChange={(e) => setNewSubName(e.target.value)}
                            onKeyDown={(e) => {
                              if (e.key === 'Enter') {
                                e.preventDefault();
                                handleAddSubcategory(p.name);
                              }
                            }}
                            className="flex-1 px-2.5 py-1 rounded-lg bg-[var(--bg-primary)] border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
                            autoFocus
                          />
                          <button
                            type="button"
                            onClick={() => handleAddSubcategory(p.name)}
                            disabled={isAddingSub || !newSubName.trim()}
                            className="px-2.5 py-1 rounded-lg bg-[var(--accent)] text-[var(--bg-primary)] font-bold text-xs transition-opacity hover:opacity-90 disabled:opacity-50 cursor-pointer"
                          >
                            {isAddingSub ? <Loader2 className="w-3 h-3 animate-spin" /> : 'Save'}
                          </button>
                          <button
                            type="button"
                            onClick={() => setAddingSubFor(null)}
                            className="px-2 py-1 rounded-lg text-xs text-[var(--fg-light)] hover:text-[var(--fg-primary)] cursor-pointer"
                          >
                            Cancel
                          </button>
                        </div>
                      )}

                    </div>
                  );
                })}
              </div>
            )}
          </div>

        </div>

        {/* Footer */}
        <div className="flex items-center justify-end px-5 py-3 border-t border-[var(--border-color)] bg-[var(--bg-darker)]/60 text-xs">
          <button
            onClick={onClose}
            className="px-4 py-1.5 rounded-lg bg-[var(--selection)] hover:bg-[var(--accent)] hover:text-[var(--bg-primary)] font-medium transition-colors cursor-pointer text-[var(--fg-primary)]"
          >
            Done
          </button>
        </div>

      </div>
    </div>
  );
};
