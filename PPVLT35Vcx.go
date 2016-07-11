;; path
(setq exec-path (append exec-path (list (expand-file-name &#34;~/test/go-projects/bin&#34;))))

;;env 

;; use emacs as default EDITOR when needed.
(setenv &#34;EDITOR&#34; &#34;emacsclient&#34;)
(server-start)
(add-hook &#39;server-visit-hook 
	  (lambda ()
	    ;; do not backup file, only in server-mode
	    (make-local-variable &#39;make-backup-files)
	    (setq make-backup-files nil)
	    (local-set-key (kbd &#34;C-c C-c&#34;) &#39;server-edit)
            ))

;; font set, nothing to do in term.

;; global setting

(add-to-list &#39;load-path &#34;~/.emacs.d&#34;)
(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp&#34;)

;; hide startup message
(setq inhibit-startup-message t)

(setq column-number-mode t)
;; use large kill-ring
(setq kill-ring-max 200)

;;(setq default-fill-column 60)

;; tab key
;;(setq-default indent-tabs-mode nil)
;;(setq default-tab-width 8)
;;(setq tab-stop-list ())
;;(loop for x downfrom 40 to 1 do
;;      (setq tab-stop-list (cons (* x 4) tab-stop-list)))

(setq sentence-end &#34;\\([。！？]\\|……\\|[.?!][]\&#34;&#39;)}]*\\($\\|[ \t]\\)\\)[ \t\n]*&#34;)
(setq sentence-end-double-space nil)

(setq enable-recursive-minibuffers t)
;;(resize-minibuffer-mode 1)



(setq default-major-mode &#39;text-mode)

(show-paren-mode t)
(setq show-paren-style &#39;parentheses)

(setq frame-title-format &#34;%b - %f&#34;)
;; default dir
(setq default-directory &#34;~&#34;)
;; C-k at line-head ,delete it
(setq-default kill-whole-line t)

;;;(global-set-key (kbd &#34;C-k&#34;) &#39;kill-line)


;; y/n =&gt; yes/no
(fset &#39;yes-or-no-p &#39;y-or-n-p)

;; show datetime
(setq display-time-24hr-format t)
(setq display-time-day-and-date t)
(display-time)

(icomplete-mode t)

(global-font-lock-mode t)

;; misc
(put &#39;set-goal-column &#39;disabled nil)
(put &#39;narrow-to-region &#39;disabled nil)
(put &#39;upcase-region &#39;disabled nil)
(put &#39;downcase-region &#39;disabled nil)
(put &#39;LaTeX-hide-environment &#39;disabled nil)


;;; Abbrev
;; ensure abbrev mode is always on
(setq-default abbrev-mode t)

;; do not bug me about saving my abbreviations
(setq save-abbrevs nil)
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.js\\&#39;&#34; . javascript-mode))
(autoload &#39;javascript-mode &#34;javascript&#34; nil t)


(mapcar
 (function (lambda (setting)
             (setq auto-mode-alist
                   (cons setting auto-mode-alist))))
 &#39;((&#34;\\.xml$&#34;.  sgml-mode)
   (&#34;\\\.bash&#34; . sh-mode)
   (&#34;\\.rdf$&#34;.  sgml-mode)
   (&#34;\\.session&#34; . emacs-lisp-mode)
   (&#34;\\.l$&#34; . c-mode)
   (&#34;\\.css$&#34; . css-mode)
;;   (&#34;\\.js$&#34; . java-mode)
   (&#34;\\.cfm$&#34; . html-mode)
   (&#34;gnus&#34; . emacs-lisp-mode)
   (&#34;\\.idl$&#34; . idl-mode)))



(setq user-full-name &#34;Ruan Chunping&#34;)
(setq user-mail-address &#34;ruanchunping@gmail.com&#34;)

(setq dired-recursive-copies &#39;top)
(setq dired-recursive-deletes &#39;top)


;; ibuffer
(require &#39;ibuffer)
(global-set-key (kbd &#34;C-x C-b&#34;) &#39;ibuffer)

;; ido
(require &#39;ido)
(ido-mode t)


;;
(global-set-key &#34;%&#34; &#39;match-paren)
          
(defun match-paren (arg)
  &#34;Go to the matching paren if on a paren; otherwise insert %.&#34;
  (interactive &#34;p&#34;)
  (cond ((looking-at &#34;\\s\(&#34;) (forward-list 1) (backward-char 1))
        ((looking-at &#34;\\s\)&#34;) (forward-char 1) (backward-list 1))
        (t (self-insert-command (or arg 1)))))



(require &#39;session)
(add-hook &#39;after-init-hook &#39;session-initialize)
(setq desktop-globals-to-save &#39;(desktop-missing-file-warning))

(load &#34;desktop&#34;)
;;(desktop-load-default)
;; defcustom
(setq desktop-path &#39;(&#34;~&#34; &#34;.&#34;))
(desktop-save-mode t)
(setq desktop-load-locked-desktop t)
(desktop-read)


;; [Home] key
(add-hook &#39;eshell-mode-hook
          (lambda ()
            (local-set-key [home] &#39;eshell-backward-argument)
            ))

;;;;;;;;;;MY KEY BIND;;;;;;;;;;;;;;;;;

(global-set-key (kbd &#34;C-x e&#34;) &#39;eshell)
(global-set-key [(f2)] &#39;speedbar)
(global-set-key [(f3)] &#39;eshell)
(global-set-key [(f9)] &#39;flyspell-prog-mode)



;; for Subversion 1.7&#43;
;;(require &#39;vc-svn17)

(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/psgml&#34;)
(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/css-mode&#34;)



;;(load-file &#34;~/.emacs.d/site-lisp/php-mode-improved.el&#34;)
(load-file &#34;~/.emacs.d/site-lisp/php-mode-1.5.0-nxhtml-1.94.el&#34;)
;;(load-file &#34;~/.emacs.d/site-lisp/php-mode.el.1.4.el&#34;)
(require &#39;php-mode)
(require &#39;flymake)
;;
;; configure css-mode

(autoload &#39;css-mode &#34;css-mode&#34;)
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.css\\&#39;&#34; . css-mode))
;;(setq cssm-indent-function #&#39;cssm-c-style-indenter)
(setq cssm-indent-function &#39;cssm-c-style-indenter)
(setq cssm-indent-level &#39;2)
;;
(add-hook &#39;php-mode-user-hook &#39;turn-on-font-lock)
;;
;; What files to invoke the new html-mode for?
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.inc\\&#39;&#34; . php-mode))
;;(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.phtml\\&#39;&#34; . html-mode))
;;(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.phpt\\&#39;&#34; . html-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.php[34]?\\&#39;&#34; . php-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.[sj]?html?\\&#39;&#34; . web-mode))
;;(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.jsp\\&#39;&#34; . html-mode))
;;

(setq web-mode-extra-comment-keywords &#39;(&#34;TODO&#34; &#34;NOTE&#34; &#34;WARNING&#34; &#34;IMPORTANT&#34; &#34;DISABLED&#34;))
(setq web-mode-engines-alist
      &#39;((&#34;php&#34;    . &#34;\\.phtml&#34;)
	(&#34;php&#34;    . &#34;\\.phpt&#34;)
        (&#34;blade&#34;  . &#34;\\.blade\\.&#34;))
)

(require &#39;web-mode)
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.phtml\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.phpt\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.tpl\\.php\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.jsp\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.as[cp]x\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.erb\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.mustache\\&#39;&#34; . web-mode))
(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.djhtml\\&#39;&#34; . web-mode))

(add-to-list &#39;auto-mode-alist &#39;(&#34;\\.html?\\&#39;&#34; . web-mode))

(defun my-web-tag-toggle() 
  &#34;toggle web-tag-auto-close-style&#34;
  (interactive)
  (if (= 0 web-mode-tag-auto-close-style) 
      (progn 
	(setq web-mode-tag-auto-close-style 2)
	(message &#34;Tag auto close enabled.&#34;)
	)
    (progn 
      (setq web-mode-tag-auto-close-style 0)
      (message &#34;Tag auto close disabled.&#34;)
      )))

(defun my-web-mode-hook ()
  &#34;Hooks for Web mode.&#34;
  (setq web-mode-markup-indent-offset 2)
  (gtags-auto-enabel)
  (linum-mode 1)
  (rainbow-mode 1)
  (msf-abbrev-mode 1)

  (local-set-key (kbd &#34;C-c C-c&#34;) &#39;web-mode-comment-or-uncomment) 


  (local-set-key (kbd &#34;C-c t&#34;) &#39;my-web-tag-toggle)

  ;; nav tag elements
  (local-set-key (kbd &#34;ESC &lt;up&gt;&#34;) &#39;web-mode-element-beginning)
  (local-set-key (kbd &#34;ESC &lt;down&gt;&#34;) &#39;web-mode-element-end)
  
  (local-set-key (kbd &#34;M-n&#34;) &#39;web-mode-element-next)
  (local-set-key (kbd &#34;M-p&#34;) &#39;web-mode-element-previous)


)
(add-hook &#39;web-mode-hook  &#39;my-web-mode-hook)

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; c-mode
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(define-key c-mode-base-map [(f7)] &#39;compile)
(global-set-key [(f7)] &#39;compile)
;;&#39;(compile-command &#34;make&#34;)

;;(define-key c-mode-base-map &#34;\e[4~&#34; [end])
;;(global-set-key &#34;\e[4~&#34; [end])
(global-set-key [select] &#39;end-of-line)




;; load up abbrevs for these modes
(require &#39;msf-abbrev)
;;
(global-msf-abbrev-mode t) ;; for all modes with abbrevs or
;;   ;; M-x msf-abbrev-mode RET ;; for only one buffer
;;
;;   ;; You may also want to make some bindings:
(global-set-key (kbd &#34;C-c l&#34;) &#39;msf-cmd-goto-root)
(global-set-key (kbd &#34;C-c a&#34;) &#39;msf-cmd-define)

;;(require &#39;msf-abbrev)
;;(setq msf-abbrev-verbose t) ;; optional
(setq msf-abbrev-root &#34;~/.emacs.d/mode-abbrevs&#34;)
;;(global-set-key (kbd &#34;C-c l&#34;) &#39;msf-abbrev-goto-root)
;;(global-set-key (kbd &#34;C-c a&#34;) &#39;msf-abbrev-define-new-abbrev-this-mode)
;;(msf-abbrev-load)


(require &#39;header2)

;;;;;;;;;;;;;;
;;(speedbar-toggle-show-all-files nil)
(speedbar-add-supported-extension &#34;.php&#34;) ; not necessarily required
(speedbar-add-supported-extension &#34;.phtml&#34;) ; not necessarily required
(speedbar-add-supported-extension &#34;.css&#34;) ; not necessarily required

(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/ecb-snap&#34;)


(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/geben-0.26/&#34;)
(autoload &#39;geben &#34;geben&#34; &#34;PHP Debugger on Emacs&#34; t)

(defun geben-eval-current-word () 
  &#34;Evaluate a word at where the cursor is pointing.&#34; 
  (interactive) 
  (let ((expr (current-word))) 
    (when expr 
      (geben-with-current-session session 
        (geben-dbgp-command-eval session expr))))) 

;;;;;;;;;;;;;;;
(custom-set-variables
 ;; custom-set-variables was added by Custom.
 ;; If you edit it by hand, you could mess it up, so be careful.
 ;; Your init file should contain only one such instance.
 ;; If there is more than one, they won&#39;t work right.
 &#39;(ansi-color-names-vector [&#34;#2e3436&#34; &#34;#a40000&#34; &#34;#4e9a06&#34; &#34;#c4a000&#34; &#34;#204a87&#34; &#34;#5c3566&#34; &#34;#729fcf&#34; &#34;#eeeeec&#34;])
 &#39;(compilation-auto-jump-to-first-error t)
 &#39;(compilation-window-height nil)
 &#39;(custom-enabled-themes (quote (nukq)))
 &#39;(custom-safe-themes (quote (&#34;bfa7fd19c4dbea2d9bcd8cf4657cb17e7837fbe4e3901cf8269efd31bc856108&#34; &#34;3ed3f5b22a6a41b6f407107775a906568b1f57a7b925115089cf2698d1c083cd&#34; &#34;9c883ed517fba16754d3f2254a414912ec4285920d4779a1a10027a09bdc6e9a&#34; &#34;e4553574be68e79a486f927311c7d5c7a3b956ae1c88ce278de6042d67e8c31f&#34; &#34;ce8c33a30ef408386f46c6c5e9bf1ccd83988a2644a08abe6a266561cdc8ced3&#34; &#34;5317e9da27a9275e8c41d51c93a7ffd559702b1cbaf5f62f6873512587b4298d&#34; &#34;90d0a606f7633f61efc5f7fb3d8ec56171e1b8f6ceee43a187f780c89011f14f&#34; &#34;5d1e52c8c8cd960b37c3e59025e24bb26bc4dda3f463bce2da0380c2fed51830&#34; &#34;ee5f0db7f0043a5f2d0c6f4481706eeba72e57226226a24044c0ebf074394115&#34; &#34;c8d1982e6cc30a39d7f630a1b9cd9abbfbec1ec4a177f02421028bcd59196de8&#34; &#34;e9980018b91b88e8ac645e3418607f900cfa4fdb72faf8ea250e295c0e6f84c9&#34; &#34;f8b947e6e39abb4bcfce4a84adfa6865829f6dc026660f34aab64b81958245d3&#34; &#34;fc6fc5db63d225fc48a66a8e7b9d38f15f6c8762db7e607c9543a2504905e6ac&#34; &#34;14fcd30219187087ede3d6904ed28b0c5374b7ca165364e396e17da46fb07999&#34; &#34;7639ddc879dc04a2253f841994c5866747653ef1008e9ba48163eec7a7ac274e&#34; &#34;405aece04af8920eca725ee01f970b97101c7f9083bd19712ad114d6024596d8&#34; &#34;cad3235cc6e02483a030dc9ca678882f572c2b76645734d636cf0b584801f68b&#34; &#34;59463731b3f2a0d83a512ff7291275172326357dc8f88df191034e31e051248b&#34; &#34;1f38bfc050197ef4e4ae81d12308de11db8e0fdcfd48fb8c0af49f6c443e45d5&#34; &#34;6eb481362fe441a474d2a676a79b390992f9ecd9fdb7d7573826b7459896e44d&#34; &#34;9b433d637e3e570b10842a29e01d979bafd60dff36c5ea10d27ee894d7981915&#34; &#34;3ba7dd33357006538944583ad5c4b6ecaa2a1a5b78895b163c12dd9c6c665b98&#34; &#34;64dc2c5a0b7a46dd2d75fe59236ea828ad1cd7227e8a7ae6644a7fec98dd4f69&#34; &#34;28605dabc49ecc34e63d8648c393392f30f322bd1f173ebb8dd86203d21ddda5&#34; &#34;5497720799eac433506c5b947a37c15328b19d63517d1cd6a130a4211764f52a&#34; &#34;b6944c72ef775c5fbe151f81a709657876c7bc351b59d3bdf4ff4d9d7036cc8e&#34; &#34;d01a26244039def4f0eb5aacf100e136f3bf906f67a0dbb3b8e3e2fceebc8ae9&#34; &#34;c55b70bf80a43f80d589115ae0858cfd03336d14952ae9284ecdd537fa9561d2&#34; &#34;8df89a0aefa61e95f2cde8cbc09a9af7424ee3a5e96d889b48591d4e9793c86a&#34; &#34;09c4265c5b30f3141164206d9bebbe78ee37d143dc9eaefc244a92444e7c3501&#34; &#34;435a8a3907f2cbe59b22a8c6d27bddaade2282f722f8c38a051c125e3164b5a7&#34; &#34;964422ce89a8f02d0d36aa9f69abcd4e6524f7604532c9bcefa49d48cb2bcc39&#34; &#34;b99801dc023ad67c0c94242c895f46f4a096e923e59ba19c158521df731999f4&#34; &#34;52b36e84b4746cbdc46ff70c198c80dcde5f514d3c3fdc88db6fc6ff6025ec9e&#34; &#34;21e71a457465b50b76f2210dce2ebdc3514f22844db2c6cce0b1e8a9622e71df&#34; &#34;805aaa3a952323c9205bc294770aeaca5cf93aa7a081343afebcb2518d131bbd&#34; &#34;70e7c5325c62608da1adb4573b536fdb285c272c156f9a76d14ba5306c9712d9&#34; &#34;e0598a8d8c83f81bcb6dd1e109723020ba8f3d14ecab394bf62a746e15a88fee&#34; &#34;db9d2e9ec605d46e85bf41d73b801eb29b65cdd0e82a1d26fd8b518b8b4ac3bf&#34; &#34;351489e3433c892e92960a3a050b8f1d6e959df68f45501a32d7e85ac14dbd53&#34; &#34;7d84523cc04d97d0d83bb5476d79008c8c71f8d28144535a9fd4f112fab332c2&#34; &#34;66e40926642de176afb11ac577b0ac65cb37a3e88e454ab52db040e52147db6b&#34; &#34;1b36822d3dff9cd0df94bdc7291418c3c85b0b69dea71a78a5b2826fc50e8738&#34; &#34;18de4a2bf2b23436506d613042ee3e88c2faa79e8ecc1d0c98662a03590686be&#34; &#34;86f23c8921d102d4ea1e826854e2ba204b3a5b900873ce19c0cbdc0435cce2e6&#34; &#34;48a194cba79b9e05d156506c84413ff8f417da1ca5baa2388ea3aef13c5cdf24&#34; &#34;50284ec6b91411de5880b91195e2337693a39f01fff766f8a1da6327ae61001e&#34; &#34;bb4ef154c5732b066eae98d825af6de8bcb59e190bbb848d661ff8e93f99a530&#34; &#34;5c527a13fffa7e5d50ed8f618e7169eea7ffacca9c6151bb9d678b5d116f6899&#34; &#34;539435b8be2734760eb618316c8d4048c3fa52b29670bf227c1f565fbb06fad6&#34; &#34;da124de323cc2abf4490760cbb34e4d6b9858ee07b06890a4295f02f6ea8d943&#34; &#34;6531a63b9b3c1b4df88883fb013b5baf70159c7110f09469b7f012491cbff1ff&#34; &#34;e13b5d14f2a14594ad79cbeb61d51e8a5e07547b8133fcc619ca024b15e1877a&#34; &#34;2f7e5a927bfe21013e58e44adc3706c9824014bcec096c5f4b482257fb39778f&#34; &#34;61ab2ad5f91654919e270ffe7937f367465406b75668706b8bbae5967e72b757&#34; &#34;d1c533587d66c8a7400561397e5f2cc1f00227c400a090caa29f3245c87f7315&#34; &#34;72ce3a09b3455023246f4fc1efdf41cd69e802db697ae4f87d19f757b1bbe4a0&#34; &#34;9fa61c86325922771e64650b59fa10814d99d6d89b3d3a49a1b27bd335707c4b&#34; &#34;02adc88d24dbd53d363925f5544f5f1d54d39b290f217d115bb547350f6a0645&#34; &#34;e9775aa4925cd5e38c37121dddd113b9fc9715fa6ea08bb977dbeaf567a64473&#34; &#34;98ca458a269b6494bb97c12cf301034f35f8c639c64106610e98d6642983292f&#34; &#34;19492b74f5688b7181b99795525b204f0384d16485e131e090b2b418ab716907&#34; &#34;d9d15ec98245e2beade815f217187434a81accb5619c3908f5e72f7c198e3fd6&#34; &#34;328e8a1c4f0d1ec0ad3dc500bd114cc51679fd1035f3a83871eb17b8ec3ae6b9&#34; &#34;85c761ee4132a23a215f8775a5caa9a68aac9a39831028940d5ad99a0f801e06&#34; default)))
 &#39;(ecb-layout-window-sizes (quote ((&#34;top1&#34; (0.4971751412429379 . 0.15) (0.24858757062146894 . 0.15) (0.2542372881355932 . 0.15)) (&#34;right1&#34; (0.14689265536723164 . 0.28) (0.14689265536723164 . 0.34) (0.14689265536723164 . 0.34)))))
 &#39;(ecb-options-version &#34;2.33beta2&#34;)
 &#39;(flymake-check-file-limit 819200)
 &#39;(flymake-log-level -1)
 &#39;(gnus-select-method (quote (nnimap &#34;gmail&#43;imap&#34; (nnimap-address &#34;imap.gmail.com&#34;) (nnimap-server-port 993) (nnimap-stream tls))))
 &#39;(gtags-pop-delete nil)
 &#39;(gtags-read-only t)
 &#39;(gtags-select-buffer-single t)
 &#39;(linum-format (quote dynamic))
 &#39;(message-send-mail-function (quote message-smtpmail-send-it))
 &#39;(php-manual-path &#34;~/.emacs.d/php_manual_en&#34;)
 &#39;(rainbow-ansi-colors (quote auto))
 &#39;(rainbow-ansi-colors-major-mode-list (quote (sh-mode c-mode c&#43;&#43;-mode html-mode)))
 &#39;(rainbow-html-colors (quote auto))
 &#39;(rainbow-html-colors-major-mode-list (quote (html-mode css-mode nxml-mode xml-mode web-mode)))
 &#39;(send-mail-function (quote smtpmail-send-it))
 &#39;(smtpmail-debug-info t)
 &#39;(smtpmail-debug-verb t)
 &#39;(smtpmail-default-smtp-server &#34;smtp.gmail.com&#34;)
 &#39;(smtpmail-local-domain nil)
 &#39;(smtpmail-smtp-server &#34;smtp.gmail.com&#34;)
 &#39;(smtpmail-smtp-service 587)
 &#39;(smtpmail-starttls-credentials (quote ((&#34;smtp.gmail.com&#34; 587 &#34;&#34; &#34;&#34;))))
 &#39;(speedbar-show-unknown-files t)
 &#39;(speedbar-verbosity-level 1)
 &#39;(sql-mysql-options (quote (&#34;-A&#34;)))
 &#39;(truncate-partial-width-windows nil)
 &#39;(vc-rcs-diff-switches &#34;-u&#34;)
 &#39;(web-mode-block-padding 2)
 &#39;(web-mode-enable-auto-indentation t)
 &#39;(web-mode-enable-auto-opening t)
 &#39;(web-mode-enable-auto-pairing t)
 &#39;(web-mode-enable-block-face t)
;; &#39;(web-mode-enable-part-face t)
 &#39;(web-mode-enable-comment-keywords t)
 &#39;(web-mode-enable-css-colorization t)
 &#39;(web-mode-enable-current-element-highlight t)
 &#39;(web-mode-tag-auto-close-style 2))
(custom-set-faces
 ;; custom-set-faces was added by Custom.
 ;; If you edit it by hand, you could mess it up, so be careful.
 ;; Your init file should contain only one such instance.
 ;; If there is more than one, they won&#39;t work right.
 &#39;(default ((t (:inherit nil :stipple nil :background &#34;color-236&#34; :foreground &#34;#e0e0e0&#34; :inverse-video nil :box nil :strike-through nil :overline nil :underline nil :slant normal :weight normal :height 1 :width normal :foundry &#34;default&#34; :family &#34;default&#34;))))
 &#39;(font-lock-comment-face ((t (:foreground &#34;color-246&#34; :slant italic))))
 &#39;(font-lock-constant-face ((t (:foreground &#34;color-32&#34;))))
 &#39;(font-lock-function-name-face ((t (:foreground &#34;color-33&#34;))))
 &#39;(font-lock-string-face ((t (:foreground &#34;color-213&#34;))))
 &#39;(font-lock-variable-name-face ((t (:foreground &#34;color-222&#34;))))
 &#39;(font-lock-warning-face ((t (:foreground &#34;brightred&#34;))))
 &#39;(highlight ((t (:background &#34;color-202&#34; :foreground &#34;#ffffff&#34;))))
 &#39;(linum ((t (:background &#34;color-236&#34; :foreground &#34;color-94&#34;))))
 &#39;(web-mode-block-face ((t (:background &#34;color-237&#34;))))
 &#39;(web-mode-comment-keyword-face ((t (:foreground &#34;#fff&#34; :weight bold))))
 &#39;(web-mode-current-element-highlight-face ((t (:background &#34;color-234&#34;))))
 &#39;(web-mode-html-tag-face ((t (:foreground &#34;color-173&#34;)))))
;; try  linum color-246


(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/emacs-w3m&#34;)

(setq browse-url-browser-function &#39;w3m-browse-url)
(autoload &#39;w3m-browse-url &#34;w3m&#34; &#34;Ask a WWW browser to show a URL.&#34; t)
;; optional keyboard short-cut
(global-set-key &#34;\C-xm&#34; &#39;browse-url-at-point)


(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/iedit&#34;)
(require &#39;iedit)
(global-set-key &#34;\C-c;&#34; &#39;iedit-mode)
(add-hook &#39;iedit-mode-hook 
	  (lambda () (
		      ;;(local-set-key (kbd &#34;M-&#39;&#34;) &#39;iedit-toggle-unmatched-lines-visible)
		      )))
(define-key iedit-mode-keymap (kbd &#34;M-&#39;&#34;) &#39;iedit-toggle-unmatched-lines-visible)


(define-key php-mode-map
  [menu-bar php php-debug]
  &#39;(&#34;PHP Debug&#34; . php-debug))

(defun php-debug ()
  (interactive)
  (shell-command
   (concat &#34;/usr/local/bin/php -l \&#34;&#34;
           (buffer-file-name)
           &#34;\&#34;&#34;)))

(define-key php-mode-map
  [menu-bar php php-run]
  &#39;(&#34;Run PHP&#34; . php-run))

(defun php-run ()
  (interactive)
  (shell-command
   (concat &#34;/usr/local/bin/php -q \&#34;&#34;
           (buffer-file-name)
           &#34;\&#34;&#34;)))


(defun toggle-php-html-mode ()
  (interactive)
  &#34;Toggle mode between PHP &amp; HTML Helper modes&#34;
  (cond ((string= mode-name &#34;HTML helper&#34;)
         (php-mode))
        ((string= mode-name &#34;PHP/lah&#34;)
         (html-mode))))

(require &#39;gtags) 


;; do not create GTAGS
(defun my-gtags-create-or-update ()
  &#34;create or update the gnu global tag file&#34;
  (interactive)
  (if (not (= 0 (call-process &#34;global&#34; nil nil nil &#34; -p&#34;))) ; tagfile doesn&#39;t exist?
      (let ((olddir default-directory)
	    )
	)
    ;;  tagfile already exists; update it

    (when gtags-mode
	(shell-command (concat &#34;global --single-update &#39;&#34; (file-name-nondirectory (buffer-file-name)) &#34;&#39; &amp;&amp; echo &#39;updated tagfile&#39;&#34;))
	)

    ))

(add-hook &#39;gtags-mode-hook 
	  (lambda()
	    (local-set-key (kbd &#34;C-c C-j&#34;) &#39;gtags-find-tag)   ; find a tag, also M-.
	    (local-set-key (kbd &#34;M-.&#34;) &#39;gtags-find-tag)   ; find a tag, also M-.
	    (local-set-key (kbd &#34;C-c M-j&#34;) &#39;gtags-find-rtag)   ; find a tag, also M-.
	    (local-set-key (kbd &#34;C-c C-k&#34;) &#39;gtags-find-symbol)   ; find a tag, also M-.
	    (local-set-key (kbd &#34;C-c M-k&#34;) &#39;gtags-find-pattern)  ; reverse tag

	    (add-hook &#39;after-save-hook &#39;my-gtags-create-or-update)
	    )
	  )  

(defun gtags-auto-enabel ()
  &#34;enable gtags automatic&#34;
  (if (not (= 0 (call-process &#34;global&#34; nil nil nil &#34; -p&#34;))) ; tagfile doesn&#39;t exist?
      (gtags-mode nil)
      (gtags-mode t))
)
(defun my-php-get-pattern ()
  (save-excursion
    (while (looking-at &#34;\\sw\\|\\s_&#34;)
      (forward-char 1))
    (if (or (re-search-backward &#34;\\sw\\|\\s_&#34;
                                (save-excursion (beginning-of-line) (point))
                                t)
            (re-search-forward &#34;\\(\\sw\\|\\s_\\)&#43;&#34;
                               (save-excursion (end-of-line) (point))
                               t))
        (progn (goto-char (match-end 0))
               (buffer-substring-no-properties
                (point)
                (progn (forward-sexp -1)
                       (while (looking-at &#34;\\s&#39;&#34;)
                         (forward-char 1))
                       (point))))
      nil)))

(defun my-php-show-def ()
  (interactive)
  (let* ((tagname (my-php-get-pattern))
	 (buf (get-buffer-create &#34;*php-def*&#34;))
	 (cmd (concat &#34;global --result ctags-x \&#34;&#34; tagname &#34;\&#34;&#34;))
         (bufstr (shell-command cmd buf))
         arglist)
    (save-excursion
      (set-buffer buf)
      (goto-char (point-min))
      (when (re-search-forward
             (format &#34;.*\\s-&#43;%s\\s-*(\\([^{]*\\))&#34; tagname)
             nil t)
        (setq arglist (buffer-substring-no-properties
                       (match-beginning 0) (match-end 0)))))
    (if arglist
        (message &#34;defination of %s:\n%s&#34; tagname arglist)
        (message &#34;Unknown def: %s\n%s&#34; tagname (shell-command-to-string cmd))
      )))


(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/php-extras&#34;)
;;(eval-after-load &#39;php-mode
(require &#39;php-extras)

;;(require &#39;popup) is loaded by ac
(defun my-php-function-show-tip ()
  &#34;show docs for symbol at point or at beginning of list if not on a symbol&#34;
  (interactive)
  (let ((s (save-excursion
             (or (symbol-at-point)
                 (progn (backward-up-list)
                        (forward-char)
                        (symbol-at-point))))))
    ;;(message &#34;:: %s&#34; (php-extras-function-documentation (my-php-get-pattern)) )

    (popup-tip (php-extras-function-documentation (my-php-get-pattern)) ;; (symbol-name s))
	       :point (point)
	       :around t
	       :scroll-bar t
	       :margin t)
    ))

(defun my-php-try-get-def ()
  &#34;try get defination&#34;
  (interactive)
  ;;(let message &#34;show-tip will:%s&#34; (php-extras-function-documentation (my-php-get-pattern)) )

  (let ((s (php-extras-function-documentation (my-php-get-pattern)) ))

    (if s
	(my-php-function-show-tip)
      (my-php-show-def)
      )))

(defun my-mark-line ()
  &#34;mark current line&#34;
  (interactive)
  (progn
    (move-beginning-of-line nil)
    (set-mark-command nil)
    (move-end-of-line nil)
    )
)
(require &#39;w3m)
(defun my-php-show-manual-local ()
  &#34;show manual in local&#34;
  (interactive)
  (let ((s (php-extras-function-documentation (my-php-get-pattern)) ))
    (if s
	(w3m-find-file (concat &#34;~/.emacs.d/php_manual_en/function.&#34; (replace-regexp-in-string &#34;_&#34; &#34;-&#34; (my-php-get-pattern)) &#34;.html&#34;))
      (message &#34;no local manual found.&#34;)))
)
(defun my-php-show-manual ()
  &#34;show manual in local&#34;
  (interactive)
  (let ((s (php-extras-function-documentation (my-php-get-pattern)) ))
    (if s
	(w3m-find-file (concat &#34;~/.emacs.d/php_manual_en/function.&#34; (replace-regexp-in-string &#34;_&#34; &#34;-&#34; (my-php-get-pattern)) &#34;.html&#34;))
      (php-search-documentation)))
)

(global-set-key (kbd &#34;C-c C-l&#34;) &#39;my-mark-line)   ; mark current line

(defun my-php-mode()
  ;; 将回车代替C-j的功能，换行的同时对齐
  ;;(define-key php-mode-map [return] &#39;newline-and-indent)
  (define-key php-mode-map [(control c) (r)] &#39;php-run)
  (define-key php-mode-map [(control c) (d)] &#39;php-debug2) ;;php-lint)

  ;;C-c C-c comment/uncomment code
  (define-key php-mode-map (kbd &#34;C-c C-c&#34;) &#39;comment-or-uncomment-region) 
  (local-set-key (kbd &#34;C-c f&#34;) &#39;flymake-mode)
 
  (local-set-key (kbd &#34;ESC &lt;up&gt;&#34;) &#39;php-beginning-of-defun) 
  (local-set-key (kbd &#34;ESC &lt;down&gt;&#34;) &#39;php-end-of-defun)
  ;;(local-set-key (kbd &#34;C-c C-d&#34;) &#39;my-php-show-def)   ; find a tag, also M-.
  (local-set-key (kbd &#34;C-c h&#34;) &#39;my-php-show-manual-local)   ; w3m-find-file in php_manual_en/
  (local-set-key (kbd &#34;C-c C-f&#34;) &#39;my-php-show-manual)   ; w3m-find-file in php_manual_en/ ,if not found. then search php.net

  (local-set-key (kbd &#34;C-c C-m&#34;) &#39;my-mark-line)   ; mark current line
  (local-set-key (kbd &#34;C-c C-l&#34;) &#39;my-mark-line)   ; mark current line

 
  (interactive)
  ;; 设置php程序的对齐风格
  (c-set-style &#34;K&amp;R&#34;)
  ;; 自动模式，在此种模式下当你键入{时，会自动根据你设置的对齐风格对齐
  (c-toggle-auto-newline 1)
  ;; 此模式下，当按Backspace时会删除最多的空格
  (c-toggle-hungry-state 1)
  ;; TAB键的宽度设置为4
  (setq c-basic-offset 4)
  ;;(setq tab-width 4)
  ;; 在菜单中加入当前Buffer的函数索引
  ;;(imenu-add-menubar-index)
  ;; 在状态条上显示当前光标在哪个函数体内部
  (which-function-mode)

  (linum-mode 1)
  (rainbow-mode 1)


  ;; php-extras
  ;; eldoc-mode
  (eldoc-mode 1)

  (local-set-key (kbd &#34;C-c C-d&#34;) &#39;my-php-try-get-def) ;;function-show-tip)

  ;;(local-set-key (kbd &#34;C-c x&#34;) #&#39;php-extras-eldoc-documentation-function)
  ;; enable gtags-mode if GTAGS exists
  (gtags-auto-enabel)
  )

(add-hook &#39;php-mode-hook &#39;my-php-mode)

;;(add-to-list &#39;load-path &#34;~/.emacs.d/elpa/eldoc-eval-0.1&#34;)
;;(require &#39;eldoc-eval)


(defun my-php-hook-function ()
  (set (make-local-variable &#39;compile-command) (format &#34;php_lint %s&#34; (buffer-file-name))))
(add-hook &#39;php-mode-hook &#39;my-php-hook-function)

(defun flymake-php-init ()
  &#34;Use php and phpcs to check the syntax and code compliance of the current file.&#34;
  (let* ((temp (flymake-init-create-temp-buffer-copy &#39;flymake-create-temp-inplace))
	 (local (file-relative-name temp (file-name-directory buffer-file-name))))
    (list &#34;php_lint&#34; (list buffer-file-name))))

;;This is the error format for : php -f somefile.php -l 
(add-to-list &#39;flymake-err-line-patterns
  &#39;(&#34;\\(Parse\\|Fatal\\) error: &#43;\\(.*?\\) in \\(.*?\\) on line \\([0-9]&#43;\\)$&#34; 3 4 nil 2))

(add-to-list &#39;flymake-allowed-file-name-masks &#39;(&#34;\\.php$&#34; flymake-php-init))
(add-hook &#39;php-mode-hook (lambda () (flymake-mode 0)))

(require &#39;compile)
(pushnew &#39;(php &#34;syntax error, \\(.*\\) in \\(.*\\) on line \\([0-9]&#43;\\)$&#34; 2 3 nil 1)
         compilation-error-regexp-alist-alist)

(setq compilation-error-regexp-alist
      (append (list &#39;php) compilation-error-regexp-alist))

(defun php-lint ()
  &#34;Performs a PHP lint-check on the current file.&#34;
  (interactive)
  ;;(setq split-height-threshold nil)
  ;;(setq split-width-threshold most-positive-fixnum)
  (compile (concat &#34;~/bin/php_lint \&#34;&#34; (buffer-file-name) &#34;\&#34;&#34;)))


(defun php-debug2 ()
  (interactive)
  (if (not (= 0 (call-process-shell-command &#34;/usr/local/bin/php&#34; nil nil nil &#34; -l &#34; (buffer-file-name) )))
    (php-lint) ;; do compilation error check
    (shell-command &#34;echo &#39;Syntax OK.&#39;&#34;)
    )
  )



(defun my-compilation-hook ()
  (when (not (get-buffer-window &#34;*compilation*&#34;))
    (save-selected-window
      (save-excursion
        (let* ((w (split-window-vertically))
               (h (window-height w)))
          (select-window w)
          (switch-to-buffer &#34;*compilation*&#34;)
          (shrink-window (- h 10)))))))
(add-hook &#39;compilation-mode-hook &#39;my-compilation-hook)

(defun my-compilation-close ()
  &#34;Run compile and resize the compile window&#34;
  (interactive)
  (progn
    (setq w (get-buffer-window &#34;*compilation*&#34;))
    (window--delete w)
    )
  )

;; C-x q close *compilation* window
(defun my-compilation-hook-close()
  (global-set-key &#34;\C-cq&#34; &#39;my-compilation-close)
)
(add-hook &#39;compilation-mode-hook &#39;my-compilation-hook-close)


;;(require &#39;nlinum)
(setq linum-format &#34;% 5d ¦ &#34;)
;;(set-window-margins nil 10);; right:2


(add-to-list &#39;load-path &#34;~/.emacs.d/elpa/rainbow-mode-0.9&#34;)
(require &#39;rainbow-mode)

(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/go-mode/&#34; t)
(require &#39;go-mode-load)

(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/goflymake&#34;)
(require &#39;go-flymake)
;;(require &#39;go-flycheck)

(require &#39;go-errcheck)

;;auto complete

(add-to-list &#39;load-path &#34;~/.emacs.d/site-lisp/auto-complete-1.3.1&#34;)
(require &#39;go-autocomplete)
(require &#39;auto-complete-config)

(add-to-list &#39;ac-dictionary-directories &#34;~/.emacs.d/site-lisp/auto-complete-1.3.1/ac-dict&#34;)

(ac-config-default)

(defun my-html-mode-hook ()
  (auto-complete-mode t))
(add-hook &#39;html-mode-hook &#39;my-html-mode-hook)
;; which mode need auto-complete

(add-to-list &#39;ac-modes &#39;html-mode)
(add-to-list &#39;ac-modes &#39;web-mode)

(defun go-compile-debug ()
  &#34;compile roject&#34;
  (interactive)
  (shell-command &#34;go build&#34;))

(defun go-compile-release ()
  &#34;clean roject&#34;
  (interactive)
  (shell-command (concat &#34;go build -ldflags -s&#34;)))

(defun go-clean-project ()
  &#34;clean roject&#34;
  (interactive)
  (shell-command &#34;go clean -i&#34;))

(defun my-revert-buffer (&amp;optional arg)
  (interactive &#34;P&#34;)
  (revert-buffer t t arg))

(defun go-my-hook ()
  ;; gocode
  (auto-complete-mode 1)
  ;;(setq ac-sources &#39;(ac-source-go))
;;  (setq default-tab-width 4)
  (setq c-basic-offset 4)
  ;;(hs-minor-mode t)
  (setq tab-width 4 indent-tabs-mode nil)
  ;;(setq indent-line-function &#39;insert-tab)
  ;;(setq tab-width 4)
  (setq show-trailing-whitespace t)

  (define-key go-mode-map (kbd &#34;C-x g d&#34;) &#39;go-compile-debug)       
  (define-key go-mode-map (kbd &#34;C-x g r&#34;) &#39;go-compile-release)   

  ;;C-c C-f 格式化代码
  (define-key go-mode-map (kbd &#34;C-c C-f&#34;) &#39;gofmt) ;;go-fmt
  (local-set-key (kbd &#34;C-c C-r&#34;) &#39;go-remove-unused-imports)
  (local-set-key (kbd &#34;C-c i&#34;) &#39;go-goto-imports)

  (local-set-key (kbd &#34;C-c e&#34;) &#39;flymake-goto-next-error)

  (local-set-key (kbd &#34;C-c C-c&#34;) &#39;comment-or-uncomment-region) 

  (add-hook &#39;before-save-hook &#39;gofmt-before-save)

  ;;;(set (make-local-variable &#39;company-backends) &#39;(company-go))
  ;;;(company-mode)

  (linum-mode 1)
  
  )

(add-hook &#39;go-mode-hook &#39;go-my-hook)



